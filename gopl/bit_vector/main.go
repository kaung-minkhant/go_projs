package main

import (
	"bytes"
	"fmt"
	"iter"
)

type IntSet struct {
	words []uint64
}

func (s *IntSet) Has(x int) bool {
	word, bit := x/64, x%64
	return word < len(s.words) && (s.words[word]&(1<<bit)) != 0
}

func (s *IntSet) Add(x int) {
	word, bit := x/64, x%64
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i >= len(s.words) {
			s.words = append(s.words, t.words[i])
		} else {
			s.words[i] |= tword
		}
	}
}

func (s *IntSet) Len() int {
	length := 0
	for _, word := range s.words {
		if word == 0 {
			continue
		}
		for i := 0; i < 64; i++ {
			if word&(1<<i) != 0 {
				length++
			}
		}
	}
	return length
}

func (s *IntSet) Remove(x int) {
	word, bit := x/64, x%64
	if word >= len(s.words) {
		return
	}
	s.words[word] &^= (1 << bit)
}

func (s *IntSet) Clear() {
	s.words = nil
}

func (s *IntSet) Copy() *IntSet {
	var c IntSet = IntSet{
		words: make([]uint64, len(s.words)),
	}
	copy(c.words, s.words)
	return &c
}
func (s *IntSet) AddAll(items ...int) {
	for _, item := range items {
		s.Add(item)
	}
}

func (s *IntSet) IntersetWith(t *IntSet) {
	for i := 0; i < len(s.words) && i < len(t.words); i++ {
		s.words[i] &= t.words[i]
	}
}

func (s *IntSet) DifferenceWith(t *IntSet) {
	for i, tword := range t.words {
		if i >= len(s.words) {
			break
		}
		s.words[i] &^= tword
	}
}

func (s *IntSet) SymmetricDifferenceWith(t *IntSet) {
  for i := 0; i < len(t.words) && i < len(s.words); i++ {
    s.words[i] ^= t.words[i]
  }
	if len(t.words) > len(s.words) {
		s.words = append(s.words, t.words[len(s.words):]...)
	}
}

func (s *IntSet) Elem() iter.Seq[int] {
  return func (yield func (int) bool ) {
    mainLoop:
    for i, word := range s.words {
      if word == 0 {
        continue
      }

      for j := 0; j < 64; j++ {
        if word & (1<<j) != 0 {
          if !yield(i*64 + j) {
            break mainLoop
          }
        }
      }
    }
  }
}

func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<j) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", i*64+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

func main() {
	var x IntSet
	x.Add(1)
	x.Add(2)
	fmt.Printf("x is %s\n", &x)
	fmt.Printf("Length of x is %d\n", x.Len())
	x.Remove(2)
	x.Remove(64)
	fmt.Printf("x is %s\n", &x)
	fmt.Printf("Length of x is %d\n", x.Len())
	x.Clear()
	fmt.Printf("x is %s\n", &x)
	fmt.Printf("Length of x is %d\n", x.Len())
	x.Clear()
	fmt.Println(x.Has(2), x.Has(3))
	var y IntSet
	y.Add(3)
	x.UnionWith(&y)

	var c *IntSet = x.Copy()
	x.AddAll(1, 7, 5)
	fmt.Printf("x is %s\n", &x)
	c.AddAll(4, 5, 6)

	fmt.Printf("c is %s\n", c)
	x.SymmetricDifferenceWith(c)
	fmt.Printf("x is %s\n", &x)
  fmt.Println("looping x:")
  for item := range x.Elem() {
    fmt.Println(item)
  }
}
