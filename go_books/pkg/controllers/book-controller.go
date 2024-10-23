package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kaung-minkhant/go_projs/go_books/pkg/models"
	"github.com/kaung-minkhant/go_projs/go_books/pkg/utils"
)

var newBook models.Book

func GetBooks(w http.ResponseWriter, r *http.Request) {
	newbooks := models.GetAllBooks()
	// json.NewEncoder(w).Encode(newbooks)
	res, _ := json.Marshal(newbooks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	bookId, err := strconv.ParseInt(params["bookId"], 0, 0)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	bookDetails, _ := models.GetBookById(bookId)
	res, _ := json.Marshal(bookDetails)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	utils.ParseBody(r, &newBook)
	book := newBook.CreateBook()
	res, _ := json.Marshal(book)
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["bookId"]
	bookId, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	book := models.DeleteBookById(bookId)
	res, _ := json.Marshal(book)
	w.WriteHeader(http.StatusAccepted)
	w.Write(res)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.ParseInt(mux.Vars(r)["bookId"], 0, 0)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var book models.Book

	utils.ParseBody(r, &book)

	existngBook, db := models.GetBookById(id)

	if book.Name != "" {
		existngBook.Name = book.Name
	}

	if book.Author != "" {
		existngBook.Author = book.Author
	}

	if book.Publication != "" {
		existngBook.Publication = book.Publication
	}

	db.Save(&existngBook)

	res, _ := json.Marshal(existngBook)
	w.WriteHeader(http.StatusAccepted)
	w.Write(res)
}
