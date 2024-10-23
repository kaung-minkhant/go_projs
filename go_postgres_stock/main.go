package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kaung-minkhant/go_projs/go_postgres_stock/router"
)

func main() {
  r := router.Router()
  fmt.Println("starting server on port 8080")

  log.Fatal(http.ListenAndServe(":8080", r))
}
