package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	mu      sync.RWMutex
	balance int
)

func deposit(amount int, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	balance += amount
	mu.Unlock()
}

func read(wg *sync.WaitGroup) {
	defer wg.Done()
	mu.RLock()
	defer mu.RUnlock()
	fmt.Println("Current balance:", balance)
	time.Sleep(100 * time.Millisecond)
}

func main() {
	var wg sync.WaitGroup
	balance = 1000

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go deposit(100, &wg)
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go read(&wg)
	}
	wg.Wait()
	fmt.Println("Final balance:", balance)
}
