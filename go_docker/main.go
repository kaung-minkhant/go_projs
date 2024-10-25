package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %q\n", html.EscapeString(r.URL.Path))
	})

  http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hi")
  })

  log.Fatal(http.ListenAndServe(":8080", nil))
}
