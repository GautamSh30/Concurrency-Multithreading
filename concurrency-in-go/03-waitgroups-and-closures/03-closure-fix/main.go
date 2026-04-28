package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	// FIX: Pass loop variable as argument to closure
	// Each goroutine gets its own copy of i
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			fmt.Println(val) // each goroutine has its own copy
		}(i)
	}
	wg.Wait()
}
