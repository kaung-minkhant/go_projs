package main

import (
	"bytes"
	"fmt"
	"strings"
)

// 12345 => 12,345
func comma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return comma(s[:n-3]) + "," + s[n-3:]
}

func comma2(s string) string {
	var buf bytes.Buffer
  dot := strings.LastIndex(s, ".")
  var toProcess string = s
  var fraction string = ""
  var sign byte
  if dot != -1 {
    toProcess = s[:dot] 
    fraction = s[dot:]
  }
  if toProcess[0] == '+' || toProcess[0] == '-' {
    sign = toProcess[0]
    toProcess = toProcess[1:]
  }
  buf.WriteByte(sign)
	reminder := len(toProcess) % 3
	buf.WriteString(toProcess[:reminder])
	for i := reminder; i < len(toProcess); i += 3 {
		buf.WriteByte(',')
		buf.WriteString(toProcess[i : i+3])
	}
  buf.WriteString(fraction)
	return buf.String()
}

func comma3(s string) string {
	var buf bytes.Buffer
  mantissaStart := 0
  if s[0] == '+' || s[0] == '-' {
    buf.WriteByte(s[0])
    mantissaStart = 1
  }
  mantissaEnd := strings.Index(s, ".")
  mantissa := s[mantissaStart: mantissaEnd]
	reminder := len(mantissa) % 3
	buf.WriteString(mantissa[:reminder])
	for i := reminder; i < len(mantissa); i += 3 {
		buf.WriteByte(',')
		buf.WriteString(mantissa[i : i+3])
	}
  buf.WriteString(s[mantissaEnd:])
	return buf.String()
}

func main() {
	fmt.Println(comma2("-1234567890.1234"))
	fmt.Println(comma2("12345.34"))
}
