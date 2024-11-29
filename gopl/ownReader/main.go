package main

import (
	"io"
	"log"

	"golang.org/x/net/html"
)

type ownReader struct {
	s string
}

func (r *ownReader) Read(p []byte) (int, error) {
	 var err error
	 var n int
	n = copy(p, r.s)
	r.s = r.s[n:]
	if len(r.s) == 0 {
		err = io.EOF
	}
	 return n, err
}

func newOwnReader(s string) *ownReader {
	return &ownReader{s: s}
}

func main() {
	input := `<html>
    <head></head>
    <body>
      <p>Hello world!<p>
    </body>
  </html>`
	_, err := html.Parse(newOwnReader(input))
	if err != nil {
		log.Fatal(err)
	}
	// reader := newOwnReader(input)
	// bytes := make([]byte, 10)
	// n, err := reader.Read(bytes)
	// for err == nil {
	// 	fmt.Println(n)
	// 	n, err = reader.Read(bytes)
	// }
	// fmt.Println(n, err)
}
