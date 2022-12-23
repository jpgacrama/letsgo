package snippetbox_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"snippetbox/cmd/server"
	"testing"
)

func TestCreateSnippet(t *testing.T) {
	addr := ":4000"
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app := &server.Application{
		Addr:     &addr,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
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
