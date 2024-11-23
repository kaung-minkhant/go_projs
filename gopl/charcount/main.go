package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	category := make(map[string]map[rune]int)
	var utflen [utf8.UTFMax + 1]int
	invalid := 0

	input := bufio.NewReader(os.Stdin)
	for {
		r, n, err := input.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Error reading %s\n", err)
			os.Exit(1)
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		for categoryName, catRange := range unicode.Properties {
			if unicode.In(r, catRange) {
				catMap := category[categoryName]
				if catMap == nil {
					catMap = make(map[rune]int)
					category[categoryName] = catMap
				}
				catMap[r]++
			}
		}
		utflen[n]++
	}

	for cat, catMap := range category {
		fmt.Printf("Category: %s\n", cat)
		for r, n := range catMap {
			fmt.Printf("%q\t%d\n", r, n)
		}
	}

	fmt.Println("len\tcount")
	for l, n := range utflen {
		if l > 0 {
			fmt.Printf("%d\t%d\n", l, n)
		}
	}
	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}
}
