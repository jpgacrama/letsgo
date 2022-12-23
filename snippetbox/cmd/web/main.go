package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"snippetbox/cmd/server"
	"snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

func parseUserInputs() (*string, *string) {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()
	return addr, dsn
}

func createLoggers() (*log.Logger, *log.Logger) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return infoLog, errorLog
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	addr, dsn := parseUserInputs()
	infoLog, errorLog := createLoggers()
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	server, err := server.CreateServer(
		&server.Application{
			Addr:     addr,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
			Snippets: &mysql.SnippetModel{DB: db},
		})
	if err == nil {
		infoLog.Printf("Starting server on %s", *addr)
		errorLog.Fatal(server.ListenAndServe())
	}
}
