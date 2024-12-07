package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Node interface{}

type CharData string

type Element struct {
	Type     xml.Name
	Attr     []xml.Attr
	Children []Node
}

func parseXMLToTree(r io.Reader) Node {
	dec := xml.NewDecoder(r)
	stack := make([]Node, 0)
	var root Node
	for {
		token, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("unexpected parsing err: %s", err)
		}

		switch tok := token.(type) {
		case xml.StartElement:
			node := &Element{
				Type:     tok.Name,
				Attr:     tok.Attr,
				Children: make([]Node, 0),
			}
			// fmt.Printf("%p\n", node)
			if len(stack) == 0 {
				// fmt.Println("root node added")
				root = node
			}
			if len(stack) != 0 {
				addToTree(stack[len(stack)-1], node)
			}
			stack = append(stack, node)
			// fmt.Println(stack)
		case xml.EndElement:
			stack = stack[:len(stack)-1]
		case xml.CharData:
      txt := strings.TrimSpace(string(tok))
			node := CharData(txt)
			if len(stack) != 0 && txt != "\n" && txt != "" {
				addToTree(stack[len(stack)-1], node)
			}
		}
	}
	return root
}

func addToTree(n, newNode Node) {
	switch v := n.(type) {
	case CharData:
		log.Fatal("cannot add children to text strings")
	case *Element:
		v.Children = append(v.Children, newNode)
	}
}

func main() {
  root := parseXMLToTree(os.Stdin)
  // fmt.Printf("root is %p\n", root)
	printTree(root)
	// fmt.Printf("%#v\n", root)
	// fmt.Printf("%#v\n", root.Children)
}

func printTree(n Node) {
	switch v := n.(type) {
	case CharData:
		fmt.Printf("The text is %q\n", v)
	case *Element:
		fmt.Printf("The node is %q\n", v.Type.Local)
		if len(v.Children) > 0 {
			fmt.Printf("Its children are: \n")
			// fmt.Printf("Its children number is: %d\n", len(v.Children))
			for _, child := range v.Children {
				// fmt.Printf("%#v\n", child)
        printTree(child)
			}
		}
	}
}
