package main

import (
	"flag"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
	"snippetbox/cmd/server"
)

func parseUserInputs() *string {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	return addr
}

func createLoggers() (*log.Logger, *log.Logger) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return infoLog, errorLog
}

func main() {
	addr := parseUserInputs()
	infoLog, errorLog := createLoggers()
	server, err := server.CreateServer(&server.Application{
		Addr:     addr,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	})
	if err == nil {
		infoLog.Printf("Starting server on %s", *addr)
		errorLog.Fatal(server.ListenAndServe())
	}
}
