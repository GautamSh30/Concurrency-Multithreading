package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func writeData(w io.Writer, data string) {
	fmt.Fprintln(w, data)
}

func main() {
	// *os.File satisfies io.Writer
	writeData(os.Stdout, "Written to stdout via io.Writer")

	// *bytes.Buffer also satisfies io.Writer
	var buf bytes.Buffer
	writeData(&buf, "Written to buffer via io.Writer")
	fmt.Println("Buffer:", buf.String())

	// Interface conceals methods not in interface
	var w io.Writer = os.Stdout
	_ = w
	// w.Close() // ERROR: io.Writer has no Close method
	// Even though *os.File has Close(), the interface hides it
}
