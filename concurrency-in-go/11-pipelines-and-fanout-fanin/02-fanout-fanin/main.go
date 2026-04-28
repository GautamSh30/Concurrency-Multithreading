package main

import (
	"fmt"
	"sync"
)

func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

func merge(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(ch <-chan int) {
		defer wg.Done()
		for n := range ch {
			out <- n
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
	in := generator(2, 3, 4, 5)

	// Fan-out: distribute work to two square goroutines
	sq1 := square(in)
	sq2 := square(in)

	// Fan-in: merge results
	for v := range merge(sq1, sq2) {
		fmt.Println(v)
	}
}
