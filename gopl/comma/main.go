package main

import "fmt"

// 12345 => 12,345
func comma(s string) string {
  n := len(s)
  if n <= 3 {
    return s
  }
  return comma(s[:n-3]) + "," + s[n-3:]
}

func main() {
  fmt.Println(comma("12345"))
  fmt.Println(comma("1234"))
}
