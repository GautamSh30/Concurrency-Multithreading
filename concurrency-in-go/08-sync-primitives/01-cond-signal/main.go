package main

import (
	"fmt"
	"sync"
)

var sharedRsc = make(map[string]string)

func main() {
	var wg sync.WaitGroup
	mu := sync.Mutex{}
	c := sync.NewCond(&mu)

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.L.Lock()
		for len(sharedRsc) == 0 {
			c.Wait()
		}
		fmt.Println("Received:", sharedRsc["rsc1"])
		c.L.Unlock()
	}()

	c.L.Lock()
	sharedRsc["rsc1"] = "data from producer"
	c.Signal()
	c.L.Unlock()

	wg.Wait()
}
