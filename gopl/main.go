package main

import "fmt"
type A struct {
  value int
}
func zeroA(a A) A {
  a.value = 0
  return a
}

func recurse(input []string, count int) {
  if count == 0 {
    return
  }
  input = append(input, "a")
  count--
  recurse(input, count)
  fmt.Println(input)
}
func main() {

  recurse(nil, 5)

  // a := A{1} 
  // zeroA(a)
  // fmt.Printf("%#v\n", a)
  // x := []int{1, 2, 3, 4, 5} // [1,2,3,4,5], len=5, cap=5
  // x = append(x, 6)          // [1,2,3,4,5,6], len=6, cap=10
  // x = append(x, 7)          // [1,2,3,4,5,6, 7], len=7, cap=10
  // a := x[4:]                // [5,6,7], len=3, cap=6
  // y := alterSlice(a)        // [10, 6, 7, 11], len=4, cap=6
  //
  // fmt.Println(x)            // [1,2,3,4,10,6,7], len=7, cap=10
  // fmt.Println(y)            // [10, 6, 7, 11]
  // fmt.Println(x[0:8])
}
// func alterSlice(a []int) []int {
//   a[0] = 10                 // [10, 6, 7], len=3, cap=6
//   a = append(a, 11)         // [10, 6, 7, 11], len=4, cap=6
//   return a
// }
