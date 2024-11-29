package main

import (
	"bufio"
	"bytes"
	"fmt"
)

type ByteCounter int

func (b *ByteCounter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	scanner.Split(bufio.ScanLines)

	count := 0
	for scanner.Scan() {
		count ++
	}

	*b += ByteCounter(count)

	if err := scanner.Err(); err != nil {
		return count, err
	}
	return count, nil
}

func main() {
	var counter ByteCounter
	fmt.Fprintf(&counter, "hello %s", "world")
	fmt.Println(counter)
}
