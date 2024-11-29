package main

import (
	"bytes"
	"fmt"
	"io"
)

type Tree struct {
	value int
	left  *Tree
	right *Tree
}

func (t *Tree) Add(n int) *Tree {
	if t == nil {
		node := &Tree{
			value: n,
		}
		return node
	}

	if n < t.value {
		t.left = t.left.Add(n)
	}
	if n > t.value {
		t.right = t.right.Add(n)
	}
	return t
}

func (t *Tree) WriteSorted(w io.Writer) {
	if t == nil {
		return
	}
	t.left.WriteSorted(w)
	fmt.Fprintf(w, "%d", t.value)
	t.right.WriteSorted(w)
}

type treePrinter struct {
	w       io.Writer
	written int
}

func (printer *treePrinter) Write(p []byte) (int, error) {
	return printer.w.Write(p)
}

func (t *Tree) String() string {
	var buf bytes.Buffer
	printer := &treePrinter{w: &buf}
	buf.WriteByte('[')
	t.WriteSorted(printer)
	buf.WriteByte(']')
	return buf.String()
}

func main() {
	items := []int{2, 1, 8, 5, 4}
	var tree *Tree
	for _, item := range items {
		tree = tree.Add(item)
	}
	fmt.Println(tree)
	// tree.WriteSorted(os.Stdout)
}
