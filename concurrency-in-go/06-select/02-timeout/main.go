package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 1)

	go func() {
		time.Sleep(5 * time.Second)
		ch <- 42
	}()

	select {
	case v := <-ch:
		fmt.Println("Received:", v)
	case <-time.After(3 * time.Second):
		fmt.Println("Timeout! No response in 3 seconds")
	}
}
