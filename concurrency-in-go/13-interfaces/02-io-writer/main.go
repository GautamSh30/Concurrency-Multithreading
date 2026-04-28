package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	// os.Stdout satisfies io.Writer
	fmt.Fprintln(os.Stdout, "Written to Stdout")

	// bytes.Buffer also satisfies io.Writer
	var buf bytes.Buffer
	fmt.Fprintln(&buf, "Written to buffer")
	fmt.Println("Buffer contains:", buf.String())
}
