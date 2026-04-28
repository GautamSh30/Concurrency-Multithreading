package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	deadline := time.Now().Add(2 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// Check if deadline is set
	if d, ok := ctx.Deadline(); ok {
		fmt.Println("Deadline set to:", d.Format(time.RFC3339))
	}

	// Simulate work
	go func() {
		// Simulate a long computation
		select {
		case <-time.After(3 * time.Second):
			fmt.Println("Work completed (won't print - deadline exceeded)")
		case <-ctx.Done():
			fmt.Println("Work cancelled:", ctx.Err())
		}
	}()

	// Wait for context to expire
	<-ctx.Done()
	fmt.Println("Main:", ctx.Err())
	time.Sleep(100 * time.Millisecond) // let goroutine print
}
