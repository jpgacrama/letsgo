package main

import (
	"database/sql"
	"flag"
	"log"
	"snippetbox/cmd/server"
	"snippetbox/pkg/models"
	"snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func parseUserInputs() (*string, *string) {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()
	return addr, dsn
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
	infoLog, errorLog := server.CreateLoggers()
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
			DB: &mysql.Database{
				DB:       db,
				InfoLog:  infoLog,
				ErrorLog: errorLog},
			Snippet: &models.Snippet{
				Title:   "O snail",
				Content: "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa",
				Expires: "7",
			},
		})
	if err == nil {
		infoLog.Printf("Starting server on %s", *addr)
		errorLog.Fatal(server.ListenAndServe())
	}
}
