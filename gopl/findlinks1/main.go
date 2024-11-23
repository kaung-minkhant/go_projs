package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
		os.Exit(1)
	}

	// for _, link := range visit(nil, doc) {
	//   fmt.Println(link)
	// }

	// outline(nil, doc)

	// freq := make(map[string]int)
	// count(freq, doc)
	// fmt.Printf("%#v\n", freq)

	printText(doc)
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	// exercise 5.1
	if c := n.FirstChild; c != nil {
		links = visit(links, c)
	}

	if c := n.NextSibling; c != nil {
		links = visit(links, c)
	}

	// for c := n.FirstChild; c != nil; c = c.NextSibling {
	//   links = visit(links, c)
	// }
	return links
}

func outline(stack []string, n *html.Node) {
	if n.Type == html.ElementNode {
		stack = append(stack, n.Data)
		fmt.Println(stack)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		outline(stack, c)
	}
}

// exercise 5.2
func count(freq map[string]int, n *html.Node) {
	if n.Type == html.ElementNode {
		freq[n.Data]++
	}

	if c := n.FirstChild; c != nil {
		count(freq, c)
	}
	if c := n.NextSibling; c != nil {
		count(freq, c)
	}
}

// exercise 5.3
func includeOnlyNewLine(input string) bool {
	status := true
	for _, char := range input {
		if char != '\n' {
			status = false
		}
	}
	return status
}
func printText(n *html.Node) {
	if n.Type == html.TextNode && n.Parent.Data != "script" && n.Parent.Data != "style" {
		text := strings.TrimSpace(n.Data)
		if !includeOnlyNewLine(text) {
			fmt.Printf("%s\n", text)
		}
	}
	if c := n.FirstChild; c != nil {
		printText(c)
	}
	if c := n.NextSibling; c != nil {
		printText(c)
	}
}
