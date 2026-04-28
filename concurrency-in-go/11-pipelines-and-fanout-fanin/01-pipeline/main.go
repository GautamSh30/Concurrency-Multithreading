package main

import "fmt"

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

func main() {
	// Set up pipeline: generator → square → square
	ch := generator(2, 3, 4, 5)
	out := square(square(ch))

	for v := range out {
		fmt.Println(v)
	}
}
