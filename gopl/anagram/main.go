package main

import "fmt"

func anagram(s1, s2 string) bool {
  freq1 := freq(s1)
  freq2 := freq(s2)
  if len(freq1) != len(freq2) {
    return false
  }
  for k, v := range freq1 {
    v2, ok := freq2[k]
    if !ok || v2 != v {
      return false
    }
  }
  return true
}

func freq (s string) map[rune]int {
  freq := make(map[rune]int)
  for _, char := range s {
    freq[char] += 1
  }
  return freq
}

func main() {
  fmt.Println(anagram("aba", "baa"))
  fmt.Println(anagram("aba", "aaa"))
}
