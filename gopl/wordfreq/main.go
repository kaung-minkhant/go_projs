package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
  input := bufio.NewScanner(os.Stdin)
  freq := make(map[string]int)
  input.Split(bufio.ScanWords)

  for input.Scan() {
    freq[input.Text()]++
  }
  if err := input.Err(); err != nil {
    fmt.Fprintf(os.Stderr, "Scan err: %s\n", err)
    os.Exit(1)
  }

  for k, v := range freq {
    fmt.Printf("%-30s: %d\n", k, v)
  }
}
