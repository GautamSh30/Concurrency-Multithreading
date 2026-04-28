package main

import (
	"bytes"
	"fmt"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		fmt.Println("Allocating new buffer")
		return new(bytes.Buffer)
	},
}

func log(msg string) {
	b := bufPool.Get().(*bytes.Buffer)
	b.Reset()
	b.WriteString(msg)
	fmt.Println(b.String())
	bufPool.Put(b)
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			log(fmt.Sprintf("Message from goroutine %d", id))
		}(i)
	}
	wg.Wait()
}
