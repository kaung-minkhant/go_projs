package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func extract(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	links := make([]string, 0)
	getLink := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					value := attr.Val
					link, err := resp.Request.URL.Parse(value)
					if err != nil {
						continue
					}
					links = append(links, link.String())
				}
			}
		}
	}

	forEachNode(doc, getLink, nil)
	return links, nil

}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}

func breathFirst(elements []string, f func(element string) []string) {
	seen := make(map[string]struct{})
	for len(elements) > 0 {
		worklist := elements
		elements = nil
		for _, item := range worklist {
			// check dup
			if _, ok := seen[item]; !ok {
				fmt.Println(item)
				seen[item] = struct{}{}
				elements = append(elements, f(item)...)
			}
		}
		fmt.Println("Length: ", len(elements))
	}
}

func crawl(url string) []string {
	links, err := extract(url)
	if err != nil {
		log.Print(err)
	}
	return links
}

func main() {
	bases := os.Args[1:]
	breathFirst(bases, crawl)
}
