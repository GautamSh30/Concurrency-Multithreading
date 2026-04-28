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

	// Two consumers waiting for data
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			c.L.Lock()
			for len(sharedRsc) == 0 {
				c.Wait()
			}
			fmt.Printf("Consumer %d received: %s\n", id, sharedRsc["rsc1"])
			c.L.Unlock()
		}(i)
	}

	// Producer writes data and wakes all consumers
	c.L.Lock()
	sharedRsc["rsc1"] = "broadcast data"
	c.Broadcast()
	c.L.Unlock()

	wg.Wait()
}
