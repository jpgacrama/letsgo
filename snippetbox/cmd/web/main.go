package main

import (
	"log"
	"net/http"
	"snippetbox/cmd/server"
)

func main() {
	log.Println("Starting server on :4000")
	server, err := server.CreateServer()
	if err == nil {
		log.Fatal(http.ListenAndServe(":4000", server.GetMux()))
	}
}
