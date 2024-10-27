package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kaung-minkhant/go_projs/go_simple_blockchain_book/blockchain"
	"github.com/kaung-minkhant/go_projs/go_simple_blockchain_book/utils"
)

var bc *blockchain.Blockchain

func main() {

  bc = blockchain.NewBlockChain()

  go func() {
    for _, block := range bc.Blocks {
      fmt.Printf("Prev Hash: %x\n", block.PrevHash)
      bytes, _ := json.MarshalIndent(block.Data, "", "  ")
      fmt.Printf("Data: %v\n", string(bytes))
      fmt.Printf("hash: %x\n", block.Hash)
      fmt.Println()
    }
  }()

  r := mux.NewRouter()

  r.HandleFunc("/", handleGetBlockchain).Methods("GET")
  r.HandleFunc("/book", handleCreateBook).Methods("POST")
  r.HandleFunc("/checkout", handleCheckout).Methods("POST")

  log.Println("starting server on port 8080")
  log.Fatal(http.ListenAndServe(":8080", r))
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
  utils.RespondJson(w, bc.Blocks, http.StatusOK)
}

func handleCreateBook(w http.ResponseWriter, r *http.Request) {
  var book blockchain.Book

  if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
    utils.RespondBadRequest(w, err)    
    return
  }

  hash := md5.New()
  io.WriteString(hash, book.Title + book.Author)

  book.ID = fmt.Sprintf("%x", hash.Sum(nil))

  utils.RespondJson(w, book, http.StatusCreated)
}

func handleCheckout(w http.ResponseWriter, r *http.Request) {
  var checkout blockchain.BookCheckout
  if err := json.NewDecoder(r.Body).Decode(&checkout); err != nil {
    utils.RespondBadRequest(w, err)
    return
  }
  
  if err := bc.WriteBlock(&checkout); err != nil {
    utils.RespondInternalError(w, err)
    return
  }

  utils.RespondJson(w, checkout, http.StatusCreated)
}
