package main

import (
	"database/sql"
	"flag"
	"snippetbox/cmd/server"
	"snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"

	"time"

	"github.com/golangcollege/sessions"

	"crypto/tls"
	"net/http"
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

// Function to set the TLS Settings
func setTLSSettings() *tls.Config {
	return &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}
}

func main() {
	flags := parseUserInputs()
	infoLog, errorLog := server.CreateLoggers()
	db, err := openDB(*flags.dsn)
	if err != nil {
		errorLog.Printf("Error Opening DB Connection: %s", err)
	}
	defer db.Close()

	templateCache, err := server.NewTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*flags.secret))
	session.Lifetime = 12 * time.Hour
	session.SameSite = http.SameSiteStrictMode
	snippets, err := mysql.NewSnippetModel(db, infoLog, errorLog)
	if err != nil {
		errorLog.Fatal(err)
	}

	tlsConfig := setTLSSettings()

	server, err := server.CreateServer(
		&server.Application{
			Port:          flags.port,
			InfoLog:       infoLog,
			ErrorLog:      errorLog,
			Snippets:      snippets,
			TemplateCache: templateCache,
			Session:       session,
			TLSConfig:     tlsConfig,
			Users:         &mysql.UserModel{DB: db}})
	if err == nil {
		infoLog.Printf("Starting server on %s", *flags.port)
		errorLog.Fatal(server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"))
	}
}
