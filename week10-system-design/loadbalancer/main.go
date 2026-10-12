package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	healthCheckInterval = 5 * time.Second
	shutdownTimeout     = 5 * time.Second
)

func main() {
	if err := run(); err != nil {
		slog.Error("load balancer stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	listenAddr := flag.String("listen", ":8080", "address the load balancer listens on")
	backendList := flag.String("backends", "", "comma-separated backend URLs")
	demo := flag.Bool("demo", false, "start 3 local echo backends on :9001-:9003, plus a dead one on :9004")
	flag.Parse()

	rawURLs := strings.Split(*backendList, ",")
	if *demo {
		rawURLs = startEchoBackends(9001, 9002, 9003)
		// Nothing listens on :9004: the health check should take it out of
		// the rotation as soon as the load balancer starts.
		rawURLs = append(rawURLs, "http://localhost:9004")
	}

	var backends []*Backend
	for _, raw := range rawURLs {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		b, err := NewBackend(raw)
		if err != nil {
			return err
		}
		backends = append(backends, b)
	}

	lb, err := NewRoundRobin(backends)
	if err != nil {
		return err
	}

	// Same graceful-shutdown shape as the investment-engine api: Ctrl+C or
	// SIGTERM cancels ctx, which stops the health checks and the server.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	lb.StartHealthChecks(ctx, healthCheckInterval)

	srv := &http.Server{
		Addr:              *listenAddr,
		Handler:           lb,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("load balancer listening", "addr", *listenAddr, "backends", len(backends))
		serverErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		return err // could not start, e.g. the port is already in use
	case <-ctx.Done():
	}

	slog.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

// startEchoBackends starts one tiny HTTP server per port that answers with
// its own port, so you can watch the rotation with curl.
func startEchoBackends(ports ...int) []string {
	urls := make([]string, 0, len(ports))
	for _, port := range ports {
		addr := fmt.Sprintf(":%d", port)
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "served by %s (client %s)\n", addr, r.Header.Get("X-Forwarded-For"))
		})
		go func() {
			if err := http.ListenAndServe(addr, mux); err != nil {
				slog.Error("echo backend stopped", "addr", addr, "error", err)
			}
		}()
		urls = append(urls, "http://localhost"+addr)
	}
	return urls
}
