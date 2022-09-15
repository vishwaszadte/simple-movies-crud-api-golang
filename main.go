package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"math/rand"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	Id string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
}

// slice of Movie structs
var movies []Movie

// GET "/" handler
func indexhandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]string)
	resp["message"] = "You're on the index page"
	json.NewEncoder(w).Encode(resp)
}

// GET "/movies" handler
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(movies)
}

// GET "/movies/{id}" handler - getting the movies by id
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for _, movie := range movies {
		if movie.Id == params["id"] {
			json.NewEncoder(w).Encode(movie)
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]string {
		"message": "movie not found",
	})
}

// POST "/movies" handler - creating a movie
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newMovie Movie
	json.NewDecoder(r.Body).Decode(&newMovie)

	newMovie.Id = strconv.Itoa(rand.Intn(10000000))

	movies = append(movies, newMovie)

	json.NewEncoder(w).Encode(newMovie)
}

// PUT "/movies/{id}" handler - update a movie with a given id
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	

	for index, movie := range movies {
		if movie.Id == params["id"] {
			movies = append(movies[ : index], movies[index + 1 : ]...)

			var updatedMovie Movie
			json.NewDecoder(r.Body).Decode(&updatedMovie)

			updatedMovie.Id = params["id"]
			movies = append(movies, updatedMovie)
			json.NewEncoder(w).Encode(updatedMovie)
			return
		}
	}

	json.NewEncoder(w).Encode(`"message": "something went wrong"`)
}

// DELETE "/movies/{id}" handler - delete a movie with the given id
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for index, movie := range movies {
		if movie.Id == params["id"] {
			movies = append(movies[ : index], movies[index + 1 : ]...)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}

	json.NewEncoder(w).Encode(`"message": "movie with this id does not exist"`)
}

// main function
func main() {
	// initializing the movies slice with some values
	movies = append(movies, Movie{
		Id: "1",
		Isbn: "1",
		Title: "The Godfather",
		Director: &Director {
			FirstName: "Francis",
			LastName: "Cappola",
		},
	})

	movies = append(movies, Movie{
		Id: "2",
		Isbn: "2",
		Title: "Titanic",
		Director: &Director {
			FirstName: "James",
			LastName: "Cameron",
		},
	})

	// router 
	r := mux.NewRouter()

	r.HandleFunc("/", indexhandler).Methods("GET")
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	http.Handle("/", r)
	
	fmt.Printf("Starting the server...\n")
	if err := http.ListenAndServe("localhost:8080", nil) ; err == nil {
		log.Fatal(err)
	}
}
