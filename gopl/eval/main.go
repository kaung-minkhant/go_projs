package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	input := `( max(1, 2) ) +  pow(1,2)`
	expr, err := Parse(input)
	if err != nil {
		log.Fatalf("error occured: %s", err)
	}

	// expr.Print(os.Stdout)
	fmt.Println(expr)

	env := Env{}
	if err := expr.CheckExpFromEnv(env); err != nil {
		log.Fatalf("error occured: %s", err)
	}

	result, err := expr.Eval(env)
	if err != nil {
		log.Fatalf("error occured: %s", err)
	}
	fmt.Printf("Result is %g\n", result)

	var w io.Writer
	w = os.Stdout
	rw := w.(io.ReadWriter)
  rw.Write([]byte("hello"))
}
