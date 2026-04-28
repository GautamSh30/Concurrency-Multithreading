package main

import "fmt"

// owner creates, writes to, and closes the channel
func owner() <-chan int {
	ch := make(chan int, 5)
	go func() {
		defer close(ch)
		for i := 0; i < 5; i++ {
			ch <- i
		}
	}()
	return ch // returns read-only channel
}

func main() {
	// consumer only has read access
	consumer := owner()
	for v := range consumer {
		fmt.Println("Received:", v)
	}
	fmt.Println("Done receiving!")
}
