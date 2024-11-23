package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

type HTMLPrettyPrinter struct {
	w                io.Writer
	depth            int
	ignoredTags      []string
	err              error
	partialTraversal bool
	idFound          bool
	idToSearch       string
}

func NewHTMLPrettyPrinter(w io.Writer) *HTMLPrettyPrinter {
	return &HTMLPrettyPrinter{
		w: w,
	}
}

func (printer *HTMLPrettyPrinter) Print(n *html.Node) error {
	printer.forEachNode(n, printer.start, printer.end)
	if err := printer.Err(); err != nil {
		return err
	}
	return nil
}

func (printer *HTMLPrettyPrinter) ElementByID(n *html.Node, id string) error {
	printer.partialTraversal = true
	printer.idToSearch = id
	printer.forEachNode(n, printer.start, printer.end)
	if err := printer.Err(); err != nil {
		return err
	}
	return nil
}

func (printer *HTMLPrettyPrinter) Err() error {
	return printer.err
}

func (printer *HTMLPrettyPrinter) forEachNode(n *html.Node, start, end func(n *html.Node)) {
	if slices.Contains(printer.ignoredTags, n.Data) {
		return
	}
	if start != nil {
		start(n)
	}
	if printer.Err() != nil {
		return
	}

	if printer.partialTraversal == false || (printer.partialTraversal == true && !printer.idFound) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			printer.forEachNode(c, start, end)
		}
	}

	if end != nil {
		end(n)
	}
	if printer.Err() != nil {
		return
	}
}

func (printer *HTMLPrettyPrinter) AddIgnoreTags(tags []string) {
	printer.ignoredTags = tags
}

func (printer *HTMLPrettyPrinter) start(n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		// element
		printer.startElement(n)
	case html.TextNode:
		// text
		printer.startText(n)
	case html.CommentNode:
		printer.startComment(n)
	}
}

func (printer *HTMLPrettyPrinter) end(n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		//element
		printer.endElement(n)
	}
}

func (printer *HTMLPrettyPrinter) startElement(n *html.Node) {
	end := ">"
	if n.FirstChild == nil {
		end = " />"
	}
	var attrs string
	var sep string = " "
	for _, attr := range n.Attr {
		if printer.partialTraversal {
			if attr.Key == "id" && attr.Val == printer.idToSearch {
        printer.idFound = true
      }
		}
		attrs += sep + fmt.Sprintf("%s=%q", attr.Key, attr.Val)
	}
	printer.depth++
	printer.printf("%*s<%s%s%s\n", 2*printer.depth, "", n.Data, attrs, end)
}

func (printer *HTMLPrettyPrinter) endElement(n *html.Node) {
	if n.FirstChild != nil {
		printer.printf("%*s</%s>\n", 2*printer.depth, "", n.Data)
	}
	printer.depth--
}

func (printer *HTMLPrettyPrinter) startText(n *html.Node) {
	text := strings.TrimSpace(n.Data)
	if len(text) == 0 {
		return
	}
	printer.printf("%*s%s\n", 2*(printer.depth+1), "", text)
}

func (printer *HTMLPrettyPrinter) startComment(n *html.Node) {
	printer.printf("<!--%s-->\n", n.Data)
}

func (printer *HTMLPrettyPrinter) printf(format string, values ...interface{}) {
	_, err := fmt.Fprintf(printer.w, format, values...)
	printer.err = err
}

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	pp := NewHTMLPrettyPrinter(os.Stdout)
	pp.AddIgnoreTags([]string{"script", "noscript", "style"})
	// pp.Print(doc)
  pp.ElementByID(doc, "dropdown-description")
}
