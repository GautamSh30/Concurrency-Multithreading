package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	// BUG: All goroutines capture the same loop variable 'i'
	// By the time goroutines execute, i is likely 3 in all of them
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(i) // captures variable, not value
		}()
	}
	wg.Wait()
}
