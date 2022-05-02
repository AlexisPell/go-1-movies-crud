package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv" // strconv.Itoa

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// In memory DB
var movies []Movie

// Handlers //
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	for _, movie := range movies {
		if movie.ID == id {
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)

	movie.ID = strconv.Itoa(rand.Intn(100000000))

	movies = append(movies, movie)

	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	// Set json content type
	w.Header().Set("Content-Type", "application/json")
	// params
	params := mux.Vars(r)
	id := params["id"]

	// loop over the movies, range
	for index, movie := range movies {
		if movie.ID == id {
			// delete movie with the id
			movies = append(movies[:index], movies[index+1:]...)

			// add a new movie with the body
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = id
			// Put new movie into movies
			movies = append(movies, movie)
			// Return movie
			json.NewEncoder(w).Encode(movie)
			break
		}
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	res := map[string]interface{}{"id": id}

	for index, movie := range movies {
		if movie.ID == id {
			res["deleted"] = true
			// Remove the movie from slice
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}

	_, deleted := res["deleted"]
	if !deleted {
		res["deleted"] = false
	}

	json.NewEncoder(w).Encode(res)
}

// MAIN //
func main() {
	// Mux router
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "438228", Title: "Movie Two", Director: &Director{Firstname: "Steve", Lastname: "Smith"}})

	// Handlers
	// func(w http.ResponseWriter, r *http.Request)
	r.HandleFunc("/movies", getMovies).Methods(http.MethodGet)
	r.HandleFunc("/movies/{id}", getMovie).Methods(http.MethodGet)
	r.HandleFunc("/movies", createMovie).Methods(http.MethodPost)
	r.HandleFunc("/movies/{id}", updateMovie).Methods(http.MethodPut)
	r.HandleFunc("/movies/{id}", deleteMovie).Methods(http.MethodDelete)

	fmt.Printf("Starting server on 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
