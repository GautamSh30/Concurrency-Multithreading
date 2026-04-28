package main

import (
	"fmt"
	"sync"
)

func main() {
	var once sync.Once
	var wg sync.WaitGroup

	load := func() {
		fmt.Println("Initializing resource... (runs only once)")
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d: calling load\n", id)
			once.Do(load)
		}(i)
	}
	wg.Wait()
}
