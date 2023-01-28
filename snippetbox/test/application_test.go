package snippetbox_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"snippetbox/cmd/server"
	"snippetbox/pkg/models/mysql"
	"testing"
	"time"
)

var templateFiles = []string{
	"../ui/html/home.page.tmpl",
	"../ui/html/base.layout.tmpl",
	"../ui/html/footer.partial.tmpl",
}

var staticFolder = "../ui/static"

func TestHomePage(t *testing.T) {
	addr := ":4000"
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app := &server.Application{
		Addr:     &addr,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}
	t.Run("checking home page OK Case", func(t *testing.T) {
		server, err := server.CreateServer(app, templateFiles...)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodGet, "")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusOK)
	})
	t.Run("checking home page NOK Case", func(t *testing.T) {
		server, err := server.CreateServer(app, templateFiles...)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodGet, "123")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusNotFound)
	})
}

func TestStaticPage(t *testing.T) {
	server.StaticFolder = staticFolder
	addr := ":4000"
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app := &server.Application{
		Addr:     &addr,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}
	t.Run("checking static page OK Case", func(t *testing.T) {
		server, err := server.CreateServer(app, templateFiles...)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodGet, "static/")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusOK)
	})
	t.Run("checking static page NOK Case", func(t *testing.T) {
		server, err := server.CreateServer(app, templateFiles...)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodGet, "static/123")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusNotFound)
	})
}

func TestShowSnippet(t *testing.T) {
	addr := ":4000"
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	db, mock := NewMock()
	repo := &mysql.SnippetDatabase{
		DB:       db,
		InfoLog:  infoLog,
		ErrorLog: errorLog}
	defer func() {
		repo.Close()
	}()
	app := &server.Application{
		Addr:     &addr,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		DB:       repo,
	}
	t.Run("checking show snippet OK Case", func(t *testing.T) {
		server, err := server.CreateServer(app, templateFiles...)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodGet, "snippet?id=1")
		response := httptest.NewRecorder()

		// Adding ExpectPrepare to DB Expectations
		sampleDatabaseContent.ID = 1
		query := "SELECT ..."
		prep := mock.ExpectPrepare(query)
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", time.Now(), "1")
		prep.ExpectQuery().WithArgs(sampleDatabaseContent.ID).WillReturnRows(rows)

		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusOK)
	})
	t.Run("checking show snippet NOK Case", func(t *testing.T) {
		server, err := server.CreateServer(app, templateFiles...)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodGet, "snippet?id=0")
		response := httptest.NewRecorder()

		// Adding ExpectPrepare to DB Expectations
		sampleDatabaseContent.ID = 0
		query := "SELECT ..."
		prep := mock.ExpectPrepare(query)
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", time.Now(), "1")
		prep.ExpectQuery().WithArgs(sampleDatabaseContent.ID).WillReturnRows(rows)

		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusNotFound)
	})
}

func newRequest(requestType, str string) *http.Request {
	req := httptest.NewRequest(requestType, fmt.Sprintf("/%s", str), nil)
	return req
}

func assertStatus(t testing.TB, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	got := response.Code
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
