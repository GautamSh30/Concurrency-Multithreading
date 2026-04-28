package main

import (
	"context"
	"fmt"
)

func generator(ctx context.Context) <-chan int {
	ch := make(chan int)
	n := 1
	go func() {
		defer close(ch)
		for {
			select {
			case ch <- n:
				n++
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := generator(ctx)

	for i := 0; i < 5; i++ {
		fmt.Println(<-ch)
	}

	cancel() // signal generator to stop
	fmt.Println("Cancelled generator after 5 values")
}
