// package main
//
// import (
// 	"fmt"
// 	"os"
// 	"slices"
// 	"strings"
//
// 	"golang.org/x/net/html"
// )
//
// func main() {
// 	doc, err := html.Parse(os.Stdin)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
// 		os.Exit(1)
// 	}
//
// 	forEachNode(doc, startElement, endElement)
//
// }
//
// func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
// 	if pre != nil {
// 		pre(n)
// 	}
//
// 	for c := n.FirstChild; c != nil; c = c.NextSibling {
// 		forEachNode(c, pre, post)
// 	}
//
// 	if post != nil {
// 		post(n)
// 	}
// }
//
// var depth int
//
// var ommittedTags []string = []string{
// 	"script", "link", "noscript", "style",
// }
//
// func shouldOmitTag(n *html.Node) (omit bool) {
// 	if slices.Contains(ommittedTags, n.Data) {
// 		omit = true
// 	}
// 	return
// }
//
// func containsOnlyNewlines(input string) bool {
// 	onlyNewlines := true
// 	for _, char := range input {
// 		if char != '\n' {
// 			onlyNewlines = false
// 		}
// 	}
// 	return onlyNewlines
// }
//
// func startElement(n *html.Node) {
// 	if n.Type == html.ElementNode && !shouldOmitTag(n) {
// 		element := fmt.Sprintf("<%s", n.Data)
// 		var attributes string
// 		if len(n.Attr) != 0 {
// 			var sep string = " "
// 			for _, attr := range n.Attr {
// 				attributes += sep + fmt.Sprintf("%s=%q", attr.Key, attr.Val)
//
// 			}
// 		}
// 		element += attributes
// 		if n.FirstChild != nil {
// 			element += ">"
// 		} else {
// 			element += " />"
// 		}
// 		fmt.Printf("%*s%s\n", depth*2, "", element)
// 		if n.FirstChild != nil {
// 			depth++
// 		}
// 	}
// 	if n.Type == html.TextNode {
// 		var textDepth = depth + 1
// 		fmt.Printf("%*s%x\n", 2*textDepth, "", strings.TrimSpace(n.Data))
// 	}
// }
// func endElement(n *html.Node) {
// 	if n.Type == html.ElementNode && !shouldOmitTag(n) && n.FirstChild != nil {
// 		depth--
// 		fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
// 	}
// }
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var depth int

type PrettyPrinter struct {
	w   io.Writer
	err error
}

func NewPrettyPrinter() PrettyPrinter {
	return PrettyPrinter{}
}

func (pp PrettyPrinter) Pretty(w io.Writer, n *html.Node) error {
	pp.w = w
	pp.err = nil
	pp.forEachNode(n, pp.start, pp.end)
	return pp.Err()
}

func (pp PrettyPrinter) Err() error {
	return pp.err
}

func (pp PrettyPrinter) forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	if pp.Err() != nil {
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		pp.forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
	if pp.Err() != nil {
		return
	}
}

func (pp PrettyPrinter) printf(format string, args ...interface{}) {
	_, err := fmt.Fprintf(pp.w, format, args...)
	pp.err = err
}

func (pp PrettyPrinter) startElement(n *html.Node) {
	end := ">"
	if n.FirstChild == nil {
		end = "/>"
	}

	attrs := make([]string, 0, len(n.Attr))
	for _, a := range n.Attr {
		attrs = append(attrs, fmt.Sprintf(`%s="%s"`, a.Key, a.Val))
	}
	attrStr := ""
	if len(n.Attr) > 0 {
		attrStr = " " + strings.Join(attrs, " ")
	}

	name := n.Data

	pp.printf("%*s<%s%s%s\n", depth*2, "", name, attrStr, end)
	depth++
}

func (pp PrettyPrinter) endElement(n *html.Node) {
	depth--
	if n.FirstChild == nil {
		return
	}
	pp.printf("%*s</%s>\n", depth*2, "", n.Data)
}

func (pp PrettyPrinter) startText(n *html.Node) {
	text := strings.TrimSpace(n.Data)
	if len(text) == 0 {
		return
	}
	pp.printf("%*s%q\n", depth*2, "", n.Data)
}

func (pp PrettyPrinter) startComment(n *html.Node) {
	pp.printf("<!--%s-->\n", n.Data)
}

func (pp PrettyPrinter) start(n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		pp.startElement(n)
	case html.TextNode:
		pp.startText(n)
	case html.CommentNode:
		pp.startComment(n)
	}
}

func (pp PrettyPrinter) end(n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		pp.endElement(n)
	}
}

func main2() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	pp := NewPrettyPrinter()
	pp.Pretty(os.Stdout, doc)
}
