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

var port = ":4000"
var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

func TestHomePage(t *testing.T) {
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

	app := &server.Application{
		Port:          &port,
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
		SnippetDB:     repo,
		TemplateCache: templateCache,
	}
	t.Run("checking home page OK Case", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}

		// Adding ExpectPrepare to DB Expectations
		sampleDatabaseContent.ID = 1
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", time.Now(), "2024-01-24T10:23:42Z")
		prep.ExpectQuery().WillReturnRows(rows)

		request := newRequest(http.MethodGet, "")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusOK)
	})
	t.Run("checking home page NOK Case", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodGet, "123")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("checking home page NOK Case - POST instead of GET", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodPost, "")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusMethodNotAllowed)
	})
}

func TestStaticPage(t *testing.T) {
	server.StaticFolder = "../ui/static"
	app := &server.Application{
		Port:     &port,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}
	t.Run("checking static page OK Case", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodGet, "static/")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusOK)
	})
	t.Run("checking static page NOK Case", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodGet, "static/123")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("checking static page NOK Case - POST instead of GET", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodPost, "static/123")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusMethodNotAllowed)
	})

}

func TestShowSnippet(t *testing.T) {

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

	app := &server.Application{
		Port:          &port,
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
		SnippetDB:     repo,
		TemplateCache: templateCache,
	}
	t.Run("checking show snippet OK Case", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodGet, "snippet/1")
		response := httptest.NewRecorder()

		// Adding ExpectPrepare to DB Expectations
		sampleDatabaseContent.ID = 1
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", time.Now(), "2024-01-24T10:23:42Z")
		prep.ExpectQuery().WithArgs(sampleDatabaseContent.ID).WillReturnRows(rows)

		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusOK)
	})
	t.Run("checking show snippet NOK Case", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodGet, "snippet?id=0")
		response := httptest.NewRecorder()

		// Adding ExpectPrepare to DB Expectations
		sampleDatabaseContent.ID = 0
		mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", time.Now(), "2024-01-24T10:23:42Z")
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
