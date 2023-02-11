package main

import (
	"database/sql"
	"flag"
	"snippetbox/cmd/server"

	_ "github.com/go-sql-driver/mysql"

	"github.com/golangcollege/sessions"
	"time"
)

type flags struct {
	port   *string
	dsn    *string
	secret *string
}

func parseUserInputs() *flags {
	port, dsn, secret := new(string), new(string), new(string)
	if !flag.Parsed() {
		port = flag.String("port", ":4000", "HTTP network address")
		dsn = flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

		// TODO: Change the secret string to your choice
		secret = flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	}

	appFlags := &flags{
		port:   port,
		dsn:    dsn,
		secret: secret,
	}
	flag.Parse()
	return appFlags
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
	flags := parseUserInputs()
	infoLog, errorLog := server.CreateLoggers()
	db, err := openDB(*flags.dsn)
	if err != nil {
		errorLog.Fatalf("Error Opening DB Connection: %s", err)
	}
	defer db.Close()

	templateCache, err := server.NewTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*flags.secret))
	session.Lifetime = 12 * time.Hour

	server, err := server.CreateServer(
		&server.Application{
			Port:          flags.port,
			InfoLog:       infoLog,
			ErrorLog:      errorLog,
			TemplateCache: templateCache,
			Session:       session})
	if err == nil {
		infoLog.Printf("Starting server on %s", *flags.port)
		errorLog.Fatal(server.ListenAndServe())
	}
}
