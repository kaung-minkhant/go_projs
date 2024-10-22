package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)

	movie.ID = strconv.Itoa(rand.Intn(10000))
	movies = append(movies, movie)

	json.NewEncoder(w).Encode(movie)
}

// delete movie and replace
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var updatedMovie Movie
	_ = json.NewDecoder(r.Body).Decode(&updatedMovie)
	updatedMovie.ID = strconv.Itoa(rand.Intn(10000))
	// find target movie
	for index, item := range movies {
		if item.ID == params["id"] {
			// delete movie
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	movies = append(movies, updatedMovie)
	json.NewEncoder(w).Encode(updatedMovie)
}

func main() {

	movies = append(movies, Movie{
		ID:    "1",
		Isbn:  "1234",
		Title: "The meeting",
		Director: &Director{
			FirstName: "Kaung",
			LastName:  "Min Khant",
		},
	})

	movies = append(movies, Movie{
		ID:    "2",
		Isbn:  "4653",
		Title: "Our Love",
		Director: &Director{
			FirstName: "Shunn",
			LastName:  "Le Yee",
		},
	})

	r := mux.NewRouter()
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovieById).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Starting server at port 8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
