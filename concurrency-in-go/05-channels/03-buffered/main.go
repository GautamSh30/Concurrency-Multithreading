package main

import "fmt"

func main() {
	ch := make(chan int, 3)

	// Sender can send up to 3 values without blocking
	ch <- 1
	ch <- 2
	ch <- 3
	// ch <- 4 would block here (buffer full, no receiver)

	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
