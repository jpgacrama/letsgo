package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Starting server on :4000")
	mux := createRoutes()
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
