package main

import "fmt"

func noempty(strings []string) []string {
	i := 0
	for _, s := range strings {
		if s != "" {
			strings[i] = s
      i++
		}
	}
  return strings[:i]
}

func main() {
  s := []string{"abc", "", "abcd"}
  fmt.Println(noempty(s))
  fmt.Println(s)
}
