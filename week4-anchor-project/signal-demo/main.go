package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop() // releases the signal registration if main returns early

	fmt.Println("running, press Ctrl+C to stop")

	<-ctx.Done()
	fmt.Println("shutdown signal received, cleaning up for 5s...")

	stop() // restore default signal ßbehavior right away

	time.Sleep(5 * time.Second) // simulated cleanup
	fmt.Println("cleanup done")
}
