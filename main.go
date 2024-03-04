package main

import (
	"log"
	"net/http"
)

func basicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!\n"))
}

func main() {

	PORT := ":8080"

	http.HandleFunc("/", basicHandler)

	http.ListenAndServe(PORT, nil)
}