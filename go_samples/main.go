package main

import (
	"flag"
	"fmt"
	"html/template"
	"iter"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/kaung-minkhant/go_projs/go_samples/set"
)

func main() {
	// a := []byte{'a', 'b', 'c', 'd', 'e'}
	// // b := a[1:6]
	// b := a[:cap(a)]
	// fmt.Printf("a: %s\n", a)
	// fmt.Printf("b: %s\n", b)
	// receiverMain()
	yieldMain()

}

func yieldMain2() {
  
}

func PrintAllElementsPush[ E comparable ] (s *set.Set[E]) {
  s.Push(func (v E) bool {
    fmt.Println(v)
    return true
  })
}

func PrintAllElementsPull[E comparable] (s *set.Set[E]) {
  next, stop := s.Pull()
  defer stop()
  for v, ok := next(); ok; v, ok = next() {
    fmt.Println(v)
  }
}

// func Union[E comparable](s1, s2 *set.Set[E]) *set.Set[E] {
//   r := set.New[E]()
//
//   for v := range s1.m {
//     r.Add(v)
//   }
//   for v := range s2.m {
//     r.Add(v)
//   }
//
//   return r
// }

func yieldMain() {
	fibo := func() iter.Seq[int] {
		return func(yield func(int) bool) {
			a, b := 1, 1
			for {
				if !yield(b) {
					return
				}
				a, b = b, a+b
			}
		}
	}

  for n := range fibo() {
    if n > 10 {
      break
    }
    fmt.Println(n)
  }
}

func Clone1[S ~[]E, E any](s S) S {
	return s
}

type MySlice []string

func (s MySlice) String() string {
	return strings.Join(s, "+")
}
func PrintSorted(ms MySlice) string {
	c := Clone1(ms)
	slices.Sort(c)
	return c.String()
}
func typeDestructureMain() {

}

type A struct {
	ais string
}

func (r A) getValue() string {
	return r.ais
}

type B struct {
	*A
	bis string
}

func embedMain() {
	b := B{
		A: &A{
			ais: "A",
		},
		bis: "B",
	}
	fmt.Println(b.getValue())
}

type State int

const (
	ACTIVE State = iota
	INACTIVE
)

var stateName = map[State]string{
	ACTIVE:   "the state is active",
	INACTIVE: "the state is inactive",
}

func (s State) String() string {
	return stateName[s]
}

func enumMain() {
	fmt.Println("State is", enumFun(ACTIVE))
}

func enumFun(s State) State {
	switch s {
	case ACTIVE:
		return INACTIVE
	case INACTIVE:
		return ACTIVE
	default:
		panic(fmt.Errorf("Unknown state"))
	}
}

type geometry interface {
	area() int
	scale(int)
	parameter() int
}

type rect struct {
	width, height int
}

func (r rect) area() int {
	return r.width * r.height
}
func (r rect) scale(n int) {
	r.height *= n
	r.width *= n
}

func (r rect) parameter() int {
	return r.width*2 + r.height*2
}

// defined pointer receiver - value invokation works, pointer invokation works
// defined value receiver - value invokation works, pointer invokation works

func receiverMain() {
	r := []geometry{
		rect{1, 1}, &rect{2, 2},
	}
	for _, r := range r {
		fmt.Println("Area is", r.area())
		// fmt.Println("Parameter is", r.parameter())
	}
}

func sliceMain() {
	slice := make([]int, 0, 5)
	slice = slice[:cap(slice)]
	slice[0] = 1
	fmt.Println("Slice:", slice)
}

var templateStr = `
<html>
<head>
<title>QR Link Generator</title>
</head>
<body>
{{if .}}
<img src="https://quickchart.io/chart?cht=qr&chs=300x300&chl={{.}}" />
<br>
{{.}}
<br>
<br>
{{end}}
<form action="/" name=f method="GET">
    <input maxLength=1024 size=70 name=s value="" title="Text to QR Encode">
    <input type=submit value="Show QR" name=qr>
</form>
</body>
</html>
  `
var templ = template.Must(template.New("qr").Parse(templateStr))

func reserverMain() {

	var addr = flag.String("addr", ":1718", "http service address")

	flag.Parse()

	http.Handle("/", http.HandlerFunc(QR))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}

func QR(w http.ResponseWriter, r *http.Request) {
	templ.Execute(w, r.FormValue("s"))
}

func appendWithWriteMain() {
	bigSlice := make(MyByteSlice, 0, 5)
	bigSlice = MyByteSlice{'a', 'b', 'c'}

	fmt.Fprintf(&bigSlice, "Hello")

	fmt.Printf("Data: %v\n", bigSlice)
}

type MyByteSlice []byte

func (s *MyByteSlice) Write(data []byte) (n int, err error) {
	originalLength := len(*s)
	newLength := originalLength + len(data)
	n = len(data)

	// expand the original
	if newLength > cap(*s) {
		newSlice := make(MyByteSlice, newLength, newLength*2-1)
		copy(newSlice[:originalLength], *s)
		copy(newSlice[originalLength:newLength], data)
		*s = newSlice
		return
	}
	*s = (*s)[:newLength]
	copy((*s)[originalLength:newLength], data)
	return
}

func newEnumConstMain() {
	const (
		_          = iota
		KB float64 = 1 << (10 * iota)
		MB
	)

	fmt.Println("KB is", KB)
	fmt.Println("MB is", MB)
}

type human struct {
	Name string
	Age  int
}

func newFormatMain() {
	name := "shunn le"
	bytes := []byte{'a', 'b', 'c'}
	char := 65
	fmt.Printf("name: %#q\n", name)
	fmt.Printf("bytes: %v\n", bytes)
	fmt.Printf("bytes: %#q\n", bytes)
	fmt.Printf("char: %q\n", char)
	fmt.Printf("strings with x: %x\n", name)
	fmt.Printf("bytes with x: % x\n", bytes)
}

func newMain() {
	newMap := new(map[string]string)
	(*newMap)["a"] = "a"

}

func errorForMain() {
	for i := 0; i < 5; i++ {
		go func() {
			fmt.Println("i:", *(&i))
		}()
	}
	time.Sleep(3 * time.Second)
}

func fiboMain() {
	c := make(chan int)
	quit := make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- true
	}()

	fibo(c, quit)
}

func fibo(c chan int, quit chan bool) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Print("Quitting")
			return
		}
	}
}
