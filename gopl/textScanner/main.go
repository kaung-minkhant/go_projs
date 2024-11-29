package main

import (
	"fmt"
	"strings"
	"text/scanner"
	"unicode"
)

func main() {
	const src = `
// This is scanned code.
if a > 10 {
	someParsable = text
%var1 var2%

This line should not be included in the output.

/*
This multiline comment
should be extracted in
its entirety.
*/
}`

  const src2 = `aa	ab	ac	ad
ba	bb	bc	bd
ca	cb	cc	cd
da	db	dc	dd`

	var s scanner.Scanner
	s.Init(strings.NewReader(src))

	s.Filename = "example1"

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		fmt.Printf("%s: %s\n", s.Position, s.TokenText())
	}
  fmt.Println()

  s.Init(strings.NewReader(src))
  s.Filename = "example2"

  s.IsIdentRune = func (ch rune, i int) bool {
    return unicode.IsLetter(ch) || (unicode.IsDigit(ch) && i > 0) || (ch == '%' && i == 0)
  }

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		fmt.Printf("%s: %s\n", s.Position, s.TokenText())
	}
  fmt.Println()

  s.Init(strings.NewReader(src))
  s.Filename = "example3"
  s.Mode ^= scanner.SkipComments

  for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
    txt := s.TokenText()
    if strings.HasPrefix(txt, "//") || strings.HasPrefix(txt, "/*") {
      fmt.Printf("%s: %s\n", s.Position, txt)
    }
  }
  fmt.Println()

  s.Init(strings.NewReader(src2))
  var (
    col, row int
    tsv [4][4]string
  )
  s.Whitespace ^= 1 << '\t' | 1 << '\n'

  for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
    switch tok {
    case '\t': 
      col++
    case '\n':
      row++
      col = 0
    default:
      tsv[row][col] = s.TokenText()
    }
  }
  fmt.Println(tsv)
}
