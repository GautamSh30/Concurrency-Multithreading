package main

import "fmt"

func genMsg(ch chan<- string) {
	ch <- "hello"
	ch <- "world"
	close(ch)
}

func relayMsg(in <-chan string, out chan<- string) {
	for msg := range in {
		out <- msg
	}
	close(out)
}

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go genMsg(ch1)
	go relayMsg(ch1, ch2)

	for msg := range ch2 {
		fmt.Println(msg)
	}
}
