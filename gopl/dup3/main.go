package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	counts := make(map[string]int)

	for _, file := range os.Args[1:] {
		byteContent, err := os.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dup2 error: %v\n", err)
			continue
		}
		lines := strings.Split(string(byteContent), "\n")
		for _, line := range lines {
			counts[line]++
		}
	}

	for line, count := range counts {
		if count > 1 {
			fmt.Printf("%d\t%s\n", count, line)
		}
	}
}
