package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const PORT = ":8080"

func basicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!\n"))
}

func main() {
	router := chi.NewRouter()
	
	router.Use(middleware.Logger)

	router.Get("/hello", basicHandler)
	
	http.ListenAndServe(PORT, router)
}