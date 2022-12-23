package main

import (
	"flag"
	"log"
	"os"
	"snippetbox/cmd/server"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	server, err := server.CreateServer(addr, errorLog, infoLog)
	if err == nil {
		infoLog.Printf("Starting server on %s", *addr)
		errorLog.Fatal(server.ListenAndServe())
	}
}
