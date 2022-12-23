package main

import (
	"flag"
	"log"
	"net/http"
	"snippetbox/cmd/server"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	server, err := server.CreateServer()
	if err == nil {
		log.Printf("Starting server on %s", *addr)
		log.Fatal(http.ListenAndServe(*addr, server.GetMux()))
	}
}
