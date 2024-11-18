package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	fileNames := make(map[string]map[string]bool)


	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts, nil)
	} else {
		for _, file := range files {
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2 error: %v\n", err)
				continue
			}
			countLines(f, counts, fileNames)
			f.Close()
		}
	}

	for line, count := range counts {
		if count > 1 {
      sep, names := "", ""
      for filename := range fileNames[line] {
        names += sep + filename
        sep = " "
      }
			fmt.Printf("%d\t%s in %v\n", count, line, names)
		}
	}
}

func countLines(f *os.File, counts map[string]int, fileNames map[string]map[string]bool) {
	input := bufio.NewScanner(f)

	for input.Scan() {
		line := input.Text()
		counts[line]++
		if fileNames != nil {
      if _, ok := fileNames[line]; !ok {
        fileNames[line] = make(map[string]bool)
      }
      if !fileNames[line][f.Name()] {
        fileNames[line][f.Name()] = true
      }
		}
	}
}
