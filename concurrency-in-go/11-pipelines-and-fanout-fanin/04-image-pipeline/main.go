package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// walkFiles generates file paths for images in the input directory
func walkFiles(done <-chan struct{}, root string) <-chan string {
	paths := make(chan string)
	go func() {
		defer close(paths)
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
				select {
				case paths <- path:
				case <-done:
					return fmt.Errorf("walk cancelled")
				}
			}
			return nil
		})
	}()
	return paths
}

// processImage simulates processing (thumbnail generation)
func processImage(done <-chan struct{}, paths <-chan string) <-chan string {
	results := make(chan string)
	go func() {
		defer close(results)
		for path := range paths {
			select {
			case results <- fmt.Sprintf("processed: %s", filepath.Base(path)):
			case <-done:
				return
			}
		}
	}()
	return results
}

func main() {
	done := make(chan struct{})
	defer close(done)

	root := "./imgs"
	if _, err := os.Stat(root); os.IsNotExist(err) {
		fmt.Println("Note: Create an 'imgs' directory with image files to test")
		fmt.Println("This demo shows the pipeline structure:")
		fmt.Println("  walkFiles → processImage → save/print")
		return
	}

	// Pipeline: walkFiles → processImage → print results
	paths := walkFiles(done, root)
	results := processImage(done, paths)

	for r := range results {
		fmt.Println(r)
	}
}
