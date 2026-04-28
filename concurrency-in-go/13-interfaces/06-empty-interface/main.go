package main

import "fmt"

func describe(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Integer: %d (doubled: %d)\n", v, v*2)
	case string:
		fmt.Printf("String: %q (length: %d)\n", v, len(v))
	case bool:
		fmt.Printf("Boolean: %v\n", v)
	default:
		fmt.Printf("Unknown type: %T = %v\n", v, v)
	}
}

func main() {
	describe(42)
	describe("hello")
	describe(true)
	describe(3.14)
	describe([]int{1, 2, 3})
}
