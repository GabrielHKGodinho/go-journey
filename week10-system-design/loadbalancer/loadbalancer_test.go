package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

// namedServer starts a backend that answers with its own name on every path,
// including /healthz, so it always looks healthy while it is running.
func namedServer(t *testing.T, name string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, name)
	}))
	t.Cleanup(srv.Close)
	return srv
}

// newBalancer builds a RoundRobin over servers and exposes it on a test server.
func newBalancer(t *testing.T, servers ...*httptest.Server) (*RoundRobin, *httptest.Server) {
	t.Helper()
	var backends []*Backend
	for _, s := range servers {
		b, err := NewBackend(s.URL)
		if err != nil {
			t.Fatalf("NewBackend(%q): %v", s.URL, err)
		}
		backends = append(backends, b)
	}
	rr, err := NewRoundRobin(backends)
	if err != nil {
		t.Fatalf("NewRoundRobin: %v", err)
	}
	lb := httptest.NewServer(rr)
	t.Cleanup(lb.Close)
	return rr, lb
}

func get(t *testing.T, url string) (int, string) {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return resp.StatusCode, string(body)
}

func TestRoundRobinRotatesInOrder(t *testing.T) {
	_, lb := newBalancer(t, namedServer(t, "a"), namedServer(t, "b"), namedServer(t, "c"))

	var got []string
	for range 6 {
		_, body := get(t, lb.URL)
		got = append(got, body)
	}

	want := []string{"a", "b", "c", "a", "b", "c"}
	if !slices.Equal(got, want) {
		t.Errorf("rotation = %v, want %v", got, want)
	}
}

func TestRoundRobinReturnsBadGatewayBeforeHealthCheck(t *testing.T) {
	alive := namedServer(t, "alive")
	dead := namedServer(t, "dead")
	_, lb := newBalancer(t, alive, dead)
	dead.Close() // no health check has run yet, so the balancer still trusts it

	var statuses []int
	for range 4 {
		status, _ := get(t, lb.URL)
		statuses = append(statuses, status)
	}

	want := []int{http.StatusOK, http.StatusBadGateway, http.StatusOK, http.StatusBadGateway}
	if !slices.Equal(statuses, want) {
		t.Errorf("statuses = %v, want %v", statuses, want)
	}
}

func TestRoundRobinSkipsUnhealthyBackendAfterCheck(t *testing.T) {
	dead := namedServer(t, "b")
	rr, lb := newBalancer(t, namedServer(t, "a"), dead, namedServer(t, "c"))
	dead.Close()

	// Called directly: no ticker, no sleep, the test does not depend on time.
	rr.checkOnce(context.Background())

	var got []string
	for range 4 {
		status, body := get(t, lb.URL)
		if status != http.StatusOK {
			t.Fatalf("status = %d, want %d", status, http.StatusOK)
		}
		got = append(got, body)
	}

	// The rotation runs over the healthy backends only, so the load
	// stays evenly split between a and c.
	want := []string{"a", "c", "a", "c"}
	if !slices.Equal(got, want) {
		t.Errorf("rotation = %v, want %v", got, want)
	}
}

func TestRoundRobinSkipsBackendWhoseHealthzFails(t *testing.T) {
	// Up and answering, but its own health endpoint reports a problem:
	// err == nil, status 500. It must still leave the rotation.
	sick := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		io.WriteString(w, "sick")
	}))
	t.Cleanup(sick.Close)

	rr, lb := newBalancer(t, namedServer(t, "a"), sick)
	rr.checkOnce(context.Background())

	for range 4 {
		if _, body := get(t, lb.URL); body != "a" {
			t.Fatalf("body = %q, want %q: the sick backend is still in the rotation", body, "a")
		}
	}
}

func TestRoundRobinReturnsServiceUnavailableWhenNoBackendIsHealthy(t *testing.T) {
	only := namedServer(t, "only")
	rr, lb := newBalancer(t, only)
	only.Close()

	rr.checkOnce(context.Background())

	if status, _ := get(t, lb.URL); status != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want %d", status, http.StatusServiceUnavailable)
	}
}

func TestCheckOnceKeepsStateWhenContextIsCancelled(t *testing.T) {
	dead := namedServer(t, "dead")
	rr, _ := newBalancer(t, dead)
	dead.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // simulates the program shutting down

	rr.checkOnce(ctx)

	if !rr.backends[0].healthy.Load() {
		t.Error("backend was marked unhealthy during shutdown; want its last state kept")
	}
}

func TestNewRoundRobinRejectsEmptyList(t *testing.T) {
	if _, err := NewRoundRobin(nil); err != errNoBackends {
		t.Errorf("err = %v, want %v", err, errNoBackends)
	}
}
