package main

import (
	"fmt"
	"strings"
)

func basename(s string) string {
	slash := strings.LastIndex(s, "/")
	if slash != -1 {
		s = s[slash+1:]
	}
	dot := strings.LastIndex(s, ".")
	if dot != -1 {
		s = s[:dot]
	}
	// for i := len(s)-1; i>=0; i-- {
	//   if s[i] == '/' {
	//     s = s[i+1:]
	//     break
	//   }
	// }
	// for i := len(s)-1; i >= 0; i-- {
	//   if s[i] == '.' {
	//     s = s[:i]
	//     break
	//   }
	// }
	return s
}

func main() {
	fmt.Println(basename("abc"))
}
