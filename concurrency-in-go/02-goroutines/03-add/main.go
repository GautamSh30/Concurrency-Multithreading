package main

import (
	"concurrency-add/counting"
	"fmt"
)

func main() {
	numbers := counting.GenerateNumbers(1e7)

	fmt.Println("Sequential Add:")
	sum := counting.Add(numbers)
	fmt.Println("Sum:", sum)

	fmt.Println("\nConcurrent Add:")
	sum = counting.AddConcurrent(numbers)
	fmt.Println("Sum:", sum)
}
