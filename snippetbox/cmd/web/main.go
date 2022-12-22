package main

import (
	"log"
	"net/http"
	"snippetbox/cmd/handlers"
)

func main() {
	log.Println("Starting server on :4000")
	server, err := handlers.CreateServer()
	if err == nil {
		log.Fatal(http.ListenAndServe(":4000", server.GetHandler()))
	}
}
