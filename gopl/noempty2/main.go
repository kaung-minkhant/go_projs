package main

import "fmt"

func noempty(strings []string) []string {
  out := strings[:0]
  fmt.Printf("Cap: %d, Len: %d\n", cap(out), len(out))
  for _, s := range strings {
    if s != "" {
      out = append(out, s)
    }
  }
  return out
}

func main() {
  s := []string{"abc", "", "abcd"}
  fmt.Println(noempty(s))
  fmt.Println(s)
}
