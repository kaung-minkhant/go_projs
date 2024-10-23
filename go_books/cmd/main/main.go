package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kaung-minkhant/go_projs/go_books/pkg/routes"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterBookStoreRoutes(r)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
