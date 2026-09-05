package main

import (
	"context"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestNoLeaks(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctx, cancel := context.WithCancel(context.Background())
	notLeaky(ctx)
	cancel()

	time.Sleep(1 * time.Second)
}

// func TestLeaks(t *testing.T) {
// 	defer goleak.VerifyNone(t)

// 	leaky()
// }

func TestLeakyActuallyLeaks(t *testing.T) {
	leaky()
	time.Sleep(50 * time.Millisecond) // dá tempo da goroutine realmente começar e travar

	err := goleak.Find()
	if err == nil {
		t.Error("expected a goroutine leak, but none was found")
	}
}
