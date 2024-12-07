package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
	"text/scanner"
	"unicode"
)

func main() {
	dec := xml.NewDecoder(os.Stdin)
	var stack []*Tag
	var parsedArgs []*Tag
	for _, arg := range os.Args[1:] {
		parsedArgs = append(parsedArgs, parseInputTag(arg))
	}
	for {
		tok, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				log.Println("end of input")
				os.Exit(0)
			} else {
				log.Fatalf("unexpected error: %s", err)
			}
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			stack = append(stack, parseTag(&tok))
		case xml.EndElement:
			stack = stack[:len(stack)-1]
		case xml.CharData:
			// the stack needs to contain the input sequence in order
			if containsInOrder(stack, parsedArgs) {
				fmt.Printf("%s: %s\n", join(stack, " "), tok)
			}
		}
	}
}

func join(stack []*Tag, sep string) string {
	var buf bytes.Buffer
	for i, ele := range stack {
		if i != 0 {
			buf.WriteString(sep)
		}
		buf.WriteString(ele.Name)
	}
	return buf.String()
}

// x must contain y in order, can have extra or stuffs in between, but y must be in order
func containsInOrder(x, y []*Tag) (yes bool) {
	yes = false
	for len(x) >= len(y) {
		if len(y) == 0 {
			yes = true
			return
		}
		if y[0].Name == x[0].Name {
			if y[0].Ids != nil {
				for _, id := range y[0].Ids {
					if !slices.Contains(x[0].Ids, id) {
						yes = false
						return
					}
				}
			}
			if y[0].Classes != nil {
				for _, class := range y[0].Classes {
					if !slices.Contains(x[0].Classes, class) {
						yes = false
						return
					}
				}
			}
			y = y[1:]
		}
		x = x[1:]
	}
	return
}

type Tag struct {
	Name    string
	Ids     []string
	Classes []string
	Attr    []*Attr
}
type Attr struct {
	Key string
	Val string
}

func parseTag(tag *xml.StartElement) *Tag {
	var parsedTag = &Tag{}
	parsedTag.Name = tag.Name.Local
	for _, attr := range tag.Attr {
		if attr.Name.Local == "class" {
			parsedTag.Classes = append(parsedTag.Classes, strings.Split(attr.Value, " ")...)
		}
		if attr.Name.Local == "id" {
			parsedTag.Ids = append(parsedTag.Ids, strings.Split(attr.Value, " ")...)
		}
	}
	return parsedTag
}

// div => normal div
// div#id => div with id
// div.class => div with class
// must not start with # or .
func parseInputTag(tag string) *Tag {
	lexer := new(scanner.Scanner).Init(strings.NewReader(tag))
	lexer.Mode = scanner.ScanIdents
	lexer.IsIdentRune = func(ch rune, i int) bool {
		return unicode.IsLetter(ch) || ch == '-' || ch == '_' || unicode.IsDigit(ch) || ch == '/' || ch == ':'
	}
	parsedTag := &Tag{}
	for token := lexer.Scan(); token != scanner.EOF; token = lexer.Scan() {
		switch token {
		case '.':
			lexer.Scan()
			parsedTag.Classes = append(parsedTag.Classes, lexer.TokenText())
			// parsedTag.Attr = append(parsedTag.Attr, &Attr{
			// 	Key: "class",
			// 	Val: lexer.TokenText(),
			// })
		case '#':
			lexer.Scan()
			parsedTag.Ids = append(parsedTag.Ids, lexer.TokenText())
			// parsedTag.Attr = append(parsedTag.Attr, &Attr{
			// 	Key: "id",
			// 	Val: lexer.TokenText(),
			// })
		default:
			parsedTag.Name = lexer.TokenText()
		}

	}
  // fmt.Println(parsedTag)
	return parsedTag
}

func (t *Tag) String() string {
	var buf bytes.Buffer
	buf.WriteByte('<')
	buf.WriteString(t.Name)
	buf.WriteByte(' ')
	buf.WriteString(`class="`)
	buf.WriteString(strings.Join(t.Classes, " "))
	buf.WriteString(`" id="`)
	buf.WriteString(strings.Join(t.Ids, " "))
	buf.WriteString(`"`)
	buf.WriteByte('>')

	return buf.String()
}

// func main() {
// 	tag := parseTag("div.hello.world_2#world-2")
// 	fmt.Printf("%s\n", tag)
// }
