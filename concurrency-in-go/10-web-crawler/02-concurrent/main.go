package main

import (
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/net/html"
)

func findLinks(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil
	}

	var links []string
	var visit func(n *html.Node)
	visit = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			visit(c)
		}
	}
	visit(doc)
	return links
}

type result struct {
	url   string
	links []string
}

func main() {
	worklist := make(chan []string)
	var mu sync.Mutex
	fetched := make(map[string]bool)

	go func() { worklist <- []string{"https://go.dev"} }()

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for list := range worklist {
				for _, url := range list {
					mu.Lock()
					if fetched[url] {
						mu.Unlock()
						continue
					}
					fetched[url] = true
					mu.Unlock()

					fmt.Println("Fetching:", url)
					links := findLinks(url)
					if len(links) > 0 {
						go func() { worklist <- links }()
					}
				}
			}
		}()
	}

	wg.Wait()
}
