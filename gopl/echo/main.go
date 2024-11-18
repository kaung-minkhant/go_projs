package main

import (
	"fmt"
	"os"
)

func main() {
  s, sep := "", ""
  for _, value := range os.Args[1:] {
    s += sep + value
    sep = " "
  }  
  fmt.Println(s)
}
