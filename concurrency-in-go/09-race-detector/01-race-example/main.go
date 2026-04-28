package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("=== Race condition demo ===")
	fmt.Println("Run with: go run -race main.go")
	fmt.Println()

	counter := 0
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				counter++ // DATA RACE: unsynchronized access
			}
		}()
	}
	wg.Wait()
	fmt.Println("Counter (WRONG - has race):", counter)

	fmt.Println("\n=== Fixed version with mutex ===")
	var mu sync.Mutex
	safeCounter := 0

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mu.Lock()
				safeCounter++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Println("Counter (CORRECT):", safeCounter)
}
