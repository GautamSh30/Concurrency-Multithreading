package main

import (
	"fmt"
	"net/http"
	"golang.org/x/net/html"
)

var fetched = make(map[string]bool)

func findLinks(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching:", url, err)
		return nil
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error parsing:", url, err)
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

func crawl(url string, depth int) {
	if depth <= 0 || fetched[url] {
		return
	}
	fetched[url] = true
	fmt.Println("Fetching:", url)

	links := findLinks(url)
	for _, link := range links {
		crawl(link, depth-1)
	}
}

func main() {
	crawl("https://go.dev", 2)
}
