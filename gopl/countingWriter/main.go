package main

import (
	"fmt"
	"io"
	"os"
)

type byteCounter struct {
	w      io.Writer
	writen int64
}

func (b *byteCounter) Write(p []byte) (int, error) {
	n, err := b.w.Write(p)
	b.writen += int64(n)
	return n, err
}

func countingWriter(w io.Writer) (io.Writer, *int64) {
	counter := &byteCounter{w, 0}
	return counter, &counter.writen
}

func main() {
	newWritter, _ := countingWriter(os.Stdout)
	fmt.Fprintf(newWritter, "hello %s\n", "world")
}
