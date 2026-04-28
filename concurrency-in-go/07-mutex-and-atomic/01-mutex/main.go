package main

import (
	"fmt"
	"sync"
)

var (
	mu      sync.Mutex
	balance int
)

func deposit(amount int, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	balance += amount
	mu.Unlock()
}

func withdraw(amount int, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	defer mu.Unlock()
	balance -= amount
}

func main() {
	var wg sync.WaitGroup
	balance = 1000

	for i := 0; i < 5; i++ {
		wg.Add(2)
		go deposit(100, &wg)
		go withdraw(50, &wg)
	}
	wg.Wait()
	fmt.Println("Final balance:", balance)
}
