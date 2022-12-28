package snippetbox_test

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"snippetbox/cmd/server"
	"snippetbox/pkg/models/mysql"
	"testing"
)

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

func TestCreateSnippet(t *testing.T) {
	addr := ":4000"
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	db, _ := NewMock()
	app := &server.Application{
		Addr:     &addr,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Snippets: &mysql.SnippetModel{DB: db}, // mock DB
		SqlRecord: &server.Record{
			Title:   "O snail",
			Content: "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa",
			Expires: "7",
		},
	}

	t.Run("creating snippet OK Case", func(t *testing.T) {
		server, err := server.CreateServer(app, templateFiles...)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodPost, "snippet/create")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusOK)
	})
}
