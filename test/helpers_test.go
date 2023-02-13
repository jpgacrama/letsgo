package snippetbox_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golangcollege/sessions"
	"log"
	"net/http"
	"net/http/httptest"
	"snippetbox/cmd/server"
	"snippetbox/pkg/models/mysql"
	"testing"
	"time"
)

func TestHelpers(t *testing.T) {
	db, mock := NewMock()
	// New mocks due to NewSnippetModel() factory
	mock.ExpectBegin()
	_ = mock.ExpectPrepare("SELECT ...") // SELECT for Latest Statement
	_ = mock.ExpectPrepare("INSERT ...")
	prep := mock.ExpectPrepare("SELECT ...") // SELECT for just one of the items

	repo, err := mysql.NewSnippetModel(db, infoLog, errorLog)
	defer func() {
		if err == nil {
			repo.Close()
		}
	}()

	if err != nil {
		log.Fatalf("Creating NewSnippetModel failed")
		return
	}
	templateCache, err := server.NewTemplateCache("../ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*createSession()))
	session.Lifetime = 12 * time.Hour

	app := &server.Application{
		Port:          &port,
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
		Snippets:      repo,
		TemplateCache: templateCache,
		Session:       session,
	}
	t.Run("Helpers NOK Case - Template Cache does not exist", func(t *testing.T) {
		app.TemplateCache = nil
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}

		// Adding ExpectPrepare to DB Expectations
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", time.Now(), "2024-01-24T10:23:42Z")
		prep.ExpectQuery().WillReturnRows(rows)

		request := newRequest(http.MethodGet, "")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusInternalServerError)
	})
}
