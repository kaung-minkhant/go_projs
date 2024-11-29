package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"

	"golang.org/x/net/html"
)

func getElementsWithTags(n *html.Node, names ...string) []*html.Node {
	var nodes []*html.Node
	f := func(n *html.Node) {
		if n.Type == html.ElementNode && slices.Contains(names, n.Data) {
			nodes = append(nodes, n)
		}
	}
  forEachNode(n, f, nil)
  return nodes
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

func get(link string) ( *html.Node, error ) {
  resp, err := http.Get(link)
  if err != nil {
    return nil, fmt.Errorf("cannot make get request to %s: %s", link, err)
  }
  defer resp.Body.Close()
  doc, err := html.Parse(resp.Body)
  if err != nil {
    return nil, fmt.Errorf("cannot parse body from %s: %s", link, err)
  }
  return doc, nil
}

func main() {
  link := os.Args[1:][0]
  names := os.Args[2:]
  doc, err := get(link)
  if err != nil {
    log.Fatal(err)
  }
  
  nodes := getElementsWithTags(doc, names...)
  for _, node := range nodes {
    fmt.Println(node.Data)
  }
}
