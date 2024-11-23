package main

import (
	"bytes"
	"fmt"
)

func intsToString(values []int) string {
  var buf bytes.Buffer
  buf.WriteByte('[')
  for i, value := range values {
    if i > 0 {
      buf.WriteString(", ")
    }
    fmt.Fprintf(&buf, "%d", value)
  }
  buf.WriteByte(']')
  return buf.String()
}

func main() {
  fmt.Println(intsToString([]int{1,1,2,3,4,4,5,6}))
}
