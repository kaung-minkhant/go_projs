package main

import (
	"fmt"
	"unicode/utf8"
)

func rev(point []byte) {
  for i, j := 0, len(point)-1; i < j; {
    point[i], point[j] = point[j], point[i]
    i++
    j--
  }
}

func revUTF8(str []byte) []byte {
  for i := 0; i < len(str); {
    _, size := utf8.DecodeRune(str[i:])
    rev(str[i:i+size])
    i += size
  }
  rev(str)
  return str
}

func main() {
	s := []byte("Räksmörgås")
  fmt.Println(string(s))
  fmt.Println(string(revUTF8(s)))
}
