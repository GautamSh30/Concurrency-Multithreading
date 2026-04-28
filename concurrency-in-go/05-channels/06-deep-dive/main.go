package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	ch := make(chan int, 3)
	wg.Add(2)
	go producer(ch)
	go consumer(ch)
	wg.Wait()
}

func producer(ch chan<- int) {
	defer wg.Done()
	for _, v := range []int{1, 2, 3, 4, 5} {
		ch <- v
		fmt.Println("Sent:", v)
	}
	close(ch)
}

func consumer(ch <-chan int) {
	defer wg.Done()
	for v := range ch {
		fmt.Println("Received:", v)
	}
}
