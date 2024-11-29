package main

import (
	"fmt"
	"sort"
)

func IsParlindrome(s sort.Interface) bool {
	totalLength := s.Len()
	for i, j := 0, totalLength-1; i < j ; {
		if !s.Less(i, j) && !s.Less(j, i) {
			i++
			j--
      continue
		}
    return false
	}
  return true
}

type ParlindromeString struct {
  s string
}

func (p *ParlindromeString) Len() int {
  return len(p.s)
}

func (p *ParlindromeString) Less(i, j int) bool {
  return p.s[i] < p.s[j]
}

func (p *ParlindromeString) Swap(i, j int) {
  stringBytes := []byte(p.s)
  stringBytes[i], stringBytes[j] = stringBytes[j], stringBytes[i]
  p.s = string(stringBytes)
}

func main() {
  input := &ParlindromeString{"abab"}

  fmt.Println(IsParlindrome(input))
}
