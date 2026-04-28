package main

import (
	"fmt"
	"time"
)

func fun(s string) {
	for i := 0; i < 3; i++ {
		fmt.Println(s)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	// Direct function call as goroutine
	go fun("direct goroutine")

	// Anonymous function as goroutine
	go func(s string) {
		for i := 0; i < 3; i++ {
			fmt.Println(s)
			time.Sleep(100 * time.Millisecond)
		}
	}("anonymous goroutine")

	// Function value as goroutine
	fv := fun
	go fv("function value goroutine")

	// In production, use sync.WaitGroup instead of time.Sleep
	time.Sleep(time.Second)
	fmt.Println("main function done")
}
