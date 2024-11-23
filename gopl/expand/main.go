package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func usage() {
	fmt.Println(`
    Usage: ./expand key=value ...
    `)
}

func cleanUp(input string) string {
	if strings.HasPrefix(input, "${") {
		input = input[2:]
	} else if strings.HasPrefix(input, "$") {
		input = input[1:]
	}
	return input
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	subs := make(map[string]string)
	for _, arg := range os.Args[1:] {
		pairs := strings.SplitN(arg, "=", 2)
		key := pairs[0]
		value := pairs[1]

		key = cleanUp(key)
		subs[key] = value
	}

	missed := make([]string, 0,len(subs))
	f := func(input string) string {
		input = cleanUp(input)
		if _, ok := subs[input]; !ok {
			missed = append(missed, input)
			return ""
		}

		return subs[input]
	}

	var buf bytes.Buffer
	buf.ReadFrom(os.Stdin)
	fmt.Println(expand(buf.String(), f))

  missedKeys := strings.Join(missed, " ")
  if len(missed) != 0 {
    fmt.Printf("These keys are missed: %s\n", missedKeys)
  }
}

var pattern = regexp.MustCompile(`\$\w+|\${\w+}`)

func expand(input string, f func(input string) string) string {
	return pattern.ReplaceAllStringFunc(input, f)
}
