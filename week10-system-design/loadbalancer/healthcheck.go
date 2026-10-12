package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// healthCheckTimeout bounds each probe, so one hung backend cannot stall
// the whole round of checks.
const healthCheckTimeout = 1 * time.Second

// checkOnce probes every backend once and updates its healthy flag.
func (rr *RoundRobin) checkOnce(ctx context.Context) {
	for _, b := range rr.backends {
		status, err := probe(ctx, b)

		// A cancelled parent context (shutdown) makes every probe fail.
		// That says nothing about the backends, so keep their last state.
		if ctx.Err() != nil {
			return
		}

		healthy := err == nil && status == http.StatusOK

		// Swap stores the new value and returns the old one atomically,
		// so "did it change?" cannot race with another writer.
		wasHealthy := b.healthy.Swap(healthy)
		if wasHealthy == healthy {
			continue // no transition, nothing worth logging
		}

		if healthy {
			slog.Info("backend is healthy again", "backend", b.URL.Host)
		} else {
			slog.Warn("backend marked unhealthy", "backend", b.URL.Host, "status", status, "error", err)
		}
	}
}

// probe sends one GET /healthz to b and returns the status code.
// It lives in its own function so its defers run once per backend,
// not once at the end of the whole loop.
func probe(ctx context.Context, b *Backend) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, healthCheckTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, b.URL.JoinPath("healthz").String(), nil)
	if err != nil {
		return 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err // connection refused, timeout...: the backend is down
	}
	defer resp.Body.Close()

	// Draining the body lets the Transport reuse the keep-alive connection.
	_, _ = io.Copy(io.Discard, resp.Body)

	return resp.StatusCode, nil
}

// StartHealthChecks probes every backend once right away and then every
// interval, until ctx is cancelled. It returns immediately: the checks run
// in their own goroutine.
func (rr *RoundRobin) StartHealthChecks(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		rr.checkOnce(ctx) // don't wait a whole interval for the first result

		for {
			select {
			case <-ctx.Done():
				return // shutting down
			case <-ticker.C:
				rr.checkOnce(ctx)
			}
		}
	}()
}
