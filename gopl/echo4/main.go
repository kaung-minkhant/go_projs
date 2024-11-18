package main

import (
	"flag"
	"fmt"
	"math"
	"strings"
)

var n = flag.Bool("n", false, "omit trailing newline")
var sep = flag.String("sep", " ", "seperator")

func main() {
  flag.Parse()

  for x := 0; x < 20; x++ {
    fmt.Printf("x = %d eA = %8.4e\n", x, math.Exp(float64(x)))
  }

  fmt.Print(strings.Join(flag.Args(), *sep))
  if !*n {
    fmt.Println()
  }
}
