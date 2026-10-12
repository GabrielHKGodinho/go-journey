package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

// Backend is one upstream server the load balancer can forward requests to.
type Backend struct {
	URL     *url.URL
	proxy   *httputil.ReverseProxy
	healthy atomic.Bool // written by the health checker, read by every request
}

// NewBackend builds a Backend that forwards requests to rawURL.
func NewBackend(rawURL string) (*Backend, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("loadbalancer: parse backend URL %q: %w", rawURL, err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("loadbalancer: backend URL %q needs a scheme and a host", rawURL)
	}

	b := &Backend{URL: u}
	b.healthy.Store(true) // optimistic until the first check says otherwise
	b.proxy = &httputil.ReverseProxy{
		// Rewrite is the modern hook (Go 1.20+): it edits the outbound request
		// without touching the inbound one.
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(u)       // send to this backend, keeping path and query
			pr.SetXForwarded() // tell the backend who the real client was
		},
		// Without this, an unreachable backend still returns 502, but nothing
		// says which backend failed.
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			slog.Warn("backend unreachable", "backend", u.Host, "error", err)
			http.Error(w, "bad gateway", http.StatusBadGateway)
		},
	}
	return b, nil
}

// RoundRobin is an L7 load balancer that sends each request to the next
// backend in a fixed rotation.
type RoundRobin struct {
	backends []*Backend
	next     atomic.Uint64 // shared by every request goroutine, so it must be atomic
}

var errNoBackends = errors.New("loadbalancer: at least one backend is required")

// NewRoundRobin builds a RoundRobin over backends.
func NewRoundRobin(backends []*Backend) (*RoundRobin, error) {
	if len(backends) == 0 {
		return nil, errNoBackends
	}
	return &RoundRobin{backends: backends}, nil
}

// pick returns the next healthy backend in the rotation, or nil when no
// backend is healthy. The rotation runs over the healthy backends only, so
// the load stays evenly split while a backend is down.
func (rr *RoundRobin) pick() *Backend {
	// A small slice per request: negligible with a handful of backends.
	healthy := make([]*Backend, 0, len(rr.backends))
	for _, b := range rr.backends {
		if b.healthy.Load() {
			healthy = append(healthy, b)
		}
	}
	if len(healthy) == 0 {
		return nil
	}

	n := rr.next.Add(1) - 1 // Add returns the new value; subtract 1 to start at 0
	return healthy[n%uint64(len(healthy))]
}

// ServeHTTP makes RoundRobin an http.Handler.
func (rr *RoundRobin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b := rr.pick()
	if b == nil {
		// 503, not 502: we did not try any backend, we already know none can serve.
		http.Error(w, "no healthy backend", http.StatusServiceUnavailable)
		return
	}
	b.proxy.ServeHTTP(w, r)
}
