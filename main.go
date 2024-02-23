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
	ID       string   `json:"id"`
	Isbn     string   `json:"isbn"`
	Title    string   `json:"title"` // Corrected tag
	Director Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func deletMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	for index, item := range movies {
		if item.ID == param["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	for _, item := range movies {
		if item.ID == param["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

}
func createMovies(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
	body := make([]byte, 1024) // Read up to 1024 bytes of the body
	n, err := r.Body.Read(body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
	}
	fmt.Println(string(body[:n]))
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movies)
	movie.ID = strconv.Itoa(rand.Intn(100))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)

}

func updateMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	param := mux.Vars(r)
	for index, item := range movies {
		if item.ID == param["id"] {

			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}
func greeting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi greeting from Golang")
}
func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "43833", Title: "Movies One", Director: Director{Firstname: "John", Lastname: "Verma"}})
	movies = append(movies, Movie{ID: "2", Isbn: "43834", Title: "Movies Two", Director: Director{Firstname: "Black", Lastname: "Smith"}})

	r.HandleFunc("/greet", greeting).Methods("GET")
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movie/{id}", getMovie).Methods("GET")
	r.HandleFunc("/create/movies", createMovies).Methods("POST")
	r.HandleFunc("/movie/{id}", updateMovies).Methods("PUT")
	r.HandleFunc("/movie", deletMovies).Methods("DELETE")

	fmt.Println("server started... 8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}
