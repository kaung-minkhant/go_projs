package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/scanner"
)

type Tag struct {
	Name  xml.Name
	Attrs []xml.Attr
}

type lexer struct {
	lex   scanner.Scanner
	token rune
}

func (l *lexer) next() {
	l.token = l.lex.Scan()
}

func (l *lexer) text() string {
	return l.lex.TokenText()
}

func (l *lexer) describe() string {
	switch l.token {
	case scanner.EOF:
		return "end of file"
	case scanner.Ident:
		return l.text()
	case scanner.String:
		return fmt.Sprintf("%q", l.text())
	}
	return fmt.Sprintf("%c", l.token)
}

func NewLexer(r io.Reader) *lexer {
	var s scanner.Scanner
	s.Init(r)
	s.Mode = scanner.ScanIdents | scanner.ScanStrings
	// s.Whitespace = 0
	return &lexer{
		lex: s,
	}
}

// func eatWhiteSpace(l *lexer) {
// 	for l.token == ' ' || l.token == '\t' {
// 		l.next() // eat space
// 	}
// }

// tag
// tag[attr]
// tag[attr="value"]
// [attr]
// [attr="value"]

func parseInput(input string) []*Tag {
	l := NewLexer(strings.NewReader(input))
	l.next()
	selectors := make([]*Tag, 0)
	// parseSelector(l, selectors)
	for l.token != scanner.EOF {
		selectors = append(selectors, parseSelector(l))
	}
	if l.token != scanner.EOF {
		log.Fatalf("parsing went wrong, unexpected %q", l.describe())
	}
	return selectors
}

// tag
// tag[attr]
// tag[attr="value"]
// [attr]
// [attr="value"]
func parseSelector(l *lexer) *Tag {
	var selector *Tag = &Tag{
		Name: xml.Name{Local: ""},
	}
	// tag
	if l.token != '[' {
		if l.token != scanner.Ident {
			log.Fatalf("parsing falied: want tag name, got %q", l.describe())
		}
		selector.Name = xml.Name{Local: l.text()}
		l.next() // cosume tag
	}
	for l.token == '[' {
		l.next() // consume '['
		attr := parseAttribute(l)
		if attr != nil {
			selector.Attrs = append(selector.Attrs, *attr)
		}
		if l.token != ']' {
			log.Fatalf("parsing falied: want '[', got %q", l.describe())
		}
		l.next() // consume ']'
	}

	return selector
}

func parseAttribute(l *lexer) *xml.Attr {
	var attr *xml.Attr
	if l.token == ']' {
		return attr
	}
	attr = new(xml.Attr)
	switch l.token {
	case scanner.Ident:
		attr.Name = xml.Name{Local: l.text()}
		l.next() // consumne attribute name
	default:
		log.Fatal("parsing failed: want attribute name, got value string")
	}
	if l.token == '=' {
		l.next() // consume =
		switch l.token {
		case scanner.String:
			txt := l.text()
			if txt != "" {
				txt = strings.Trim(txt, `"`)
				attr.Value = txt
			}
			l.next() // consume attribute value
		default:
			log.Fatal("parsing failed: want value string")
		}
	}
	return attr
}

func matchSelectors(stack []xml.StartElement, selectors []*Tag) bool {
	for len(selectors) <= len(stack) {
		if len(selectors) == 0 {
			return true
		}
		if selectors[0].Name.Local == "" || (selectors[0].Name.Local != "" && selectors[0].Name.Local == stack[0].Name.Local) {
			if matchAttrs(stack[0].Attr, selectors[0].Attrs) {
				selectors = selectors[1:]
			}
		}

		stack = stack[1:]
	}
	return false
}

func matchAttrs(stackAttrs, selectedAttrs []xml.Attr) bool {
loop:
	for _, selectedAttr := range selectedAttrs {
		for _, stackAttr := range stackAttrs {
			if selectedAttr.Value == "" && selectedAttr.Name.Local == stackAttr.Name.Local {
				continue loop
			}

			if selectedAttr.Value != "" && selectedAttr.Name.Local == stackAttr.Name.Local && selectedAttr.Value == stackAttr.Value {
				continue loop
			}
		}
		return false
	}
	return true
}

const usage = `xmlselect2 tag tag[attr="hello"] [attr="hi"] [attr] ...`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	input := strings.Join(os.Args[1:], " ")

	selectors := parseInput(input)

	dec := xml.NewDecoder(os.Stdin)
	stack := make([]xml.StartElement, 0)
	for {
		token, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("xml parse error: %s", err)
		}

		switch tok := token.(type) {
		case xml.StartElement:
			stack = append(stack, tok)
		case xml.EndElement:
			stack = stack[:len(stack)-1]
		case xml.CharData:
			if matchSelectors(stack, selectors) {
				fmt.Printf("%s: %s\n", joinStack(stack, " "), tok)
			}
		}
	}
}

func joinStack(stack []xml.StartElement, sep string) string {
	var buf bytes.Buffer
	for i, item := range stack {
		if i > 0 {
			buf.WriteString(sep)
		}
		buf.WriteString(item.Name.Local)
	}
	return buf.String()
}
