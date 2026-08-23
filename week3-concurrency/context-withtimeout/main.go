package main

import (
	"context"
	"fmt"
	"time"
)

func slowOperation(ctx context.Context, duration time.Duration) error {
	select {
	case <-time.After(duration):
		fmt.Println("operation completed successfully")
		return nil
	case <-ctx.Done():
		fmt.Println("operation cancelled:", ctx.Err())
		return ctx.Err()
	}
}

//========= WITH TIMEOUT ================
// func main() {
// 	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 	defer cancel()

// 	err := slowOperation(ctx, 5*time.Second)
// 	if err != nil {
// 		fmt.Println("main received error:", err)
// 	}
// }

// ========= WITH DEADLINE ================
func main() {
	deadline := time.Now().Add(2 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	err := slowOperation(ctx, 5*time.Second)
	if err != nil {
		fmt.Println("main received error:", err)
	}
}
