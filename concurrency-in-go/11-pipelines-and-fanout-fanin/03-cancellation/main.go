package main

import (
	"fmt"
	"sync"
)

func generator(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func square(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func merge(done <-chan struct{}, channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(ch <-chan int) {
		defer wg.Done()
		for n := range ch {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(channels))
	for _, ch := range channels {
		go output(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	done := make(chan struct{})

	in := generator(done, 2, 3, 4, 5, 6, 7, 8)
	sq1 := square(done, in)
	sq2 := square(done, in)
	out := merge(done, sq1, sq2)

	// Only take the first value, then cancel everything
	fmt.Println("First result:", <-out)
	close(done) // Signal all goroutines to stop
}
