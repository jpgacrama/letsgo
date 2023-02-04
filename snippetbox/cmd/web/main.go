package main

import (
	"database/sql"
	"flag"
	"snippetbox/cmd/server"
	"snippetbox/pkg/models"
	"snippetbox/pkg/models/mysql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func parseUserInputs() (*string, *string) {
	port := flag.String("port", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()
	return port, dsn
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
	port, dsn := parseUserInputs()
	infoLog, errorLog := server.CreateLoggers()
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatalf("Error Opening DB Connection: %s", err)
	}
	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := server.NewTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	snippetModel, err := mysql.NewSnippetModel(db, infoLog, errorLog)
	if err != nil {
		errorLog.Fatal("Failed to create NewSnippetModel()")
		return
	}

	created := time.Now()
	expires := created.AddDate(0, 0, 1)
	server, err := server.CreateServer(
		&server.Application{
			Port:      port,
			InfoLog:   infoLog,
			ErrorLog:  errorLog,
			SnippetDB: snippetModel,
			Snippet: &models.Snippet{
				Title:   "O snail",
				Content: "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa",
				Created: created,
				Expires: expires,
			},
			TemplateCache: templateCache,
		})
	if err == nil {
		infoLog.Printf("Starting server on %s", *port)
		errorLog.Fatal(server.ListenAndServe())
	}
}
