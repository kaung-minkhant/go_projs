package main

import (
	"fmt"
	"strings"
)

func unique(input []string) []string{
  lastWritten := 0
  for _, s := range input {
    if input[lastWritten] == s {
      continue
    }
    lastWritten++
    input[lastWritten] = s
  }
  fmt.Printf("Len: %d, Cap: %d\n", len(input), cap(input))
  out := make([]string, lastWritten+1)
  copy(out, input[:lastWritten+1])
  fmt.Printf("Len: %d, Cap: %d\n", len(out), cap(out))
  // return input[:lastWritten+1]
  return out
}

func main() {
  initial := strings.Split("aabbcccdeffg", "")
  fmt.Println(unique(initial))
  fmt.Printf("initital: %v\n", initial)
}
