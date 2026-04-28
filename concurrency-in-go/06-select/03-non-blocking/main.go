package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		ch <- "message"
	}()

	select {
	case msg := <-ch:
		fmt.Println("Received:", msg)
	default:
		fmt.Println("No message received (non-blocking)")
	}
}
