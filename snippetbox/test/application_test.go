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
		Snippets:      repo,
		TemplateCache: templateCache,
	}
	t.Run("checking home page OK Case", func(t *testing.T) {
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
		Snippets:      repo,
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
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", time.Now(), "2024-01-24T10:23:42Z")
		prep.ExpectQuery().WithArgs(1).WillReturnRows(rows)

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
		mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", time.Now(), "2024-01-24T10:23:42Z")
		prep.ExpectQuery().WithArgs(0).WillReturnRows(rows)

		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("checking show snippet NOK Case - POST instead of GET", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodPost, "snippet/1")
		response := httptest.NewRecorder()

		// Adding ExpectPrepare to DB Expectations
		mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", time.Now(), "2024-01-24T10:23:42Z")
		prep.ExpectQuery().WithArgs(0).WillReturnRows(rows)

		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusMethodNotAllowed)
	})
	t.Run("checking show snippet NOK Case - Malformed snippet URL", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodGet, "snippet?id=0")
		response := httptest.NewRecorder()

		// Adding ExpectPrepare to DB Expectations
		mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", time.Now(), "2024-01-24T10:23:42Z")
		prep.ExpectQuery().WithArgs(0).WillReturnRows(rows)

		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusNotFound)
	})
}

func TestCreateSnippet(t *testing.T) {
	db, mock := NewMock()

	// New mocks due to NewSnippetModel() factory
	mock.ExpectBegin()
	_ = mock.ExpectPrepare("SELECT ...") // SELECT for Latest Statement
	prep := mock.ExpectPrepare("INSERT INTO snippets \\(title, content, created, expires\\) VALUES\\(\\?, \\?, UTC_TIMESTAMP\\(\\), DATE_ADD\\(UTC_TIMESTAMP\\(\\), INTERVAL \\? DAY\\)\\)")
	_ = mock.ExpectPrepare("SELECT ...") // SELECT for just one of the items

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
		Snippets:      repo,
		TemplateCache: templateCache,
	}
	t.Run("checking create snippet OK Case", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodPost, "snippet/create")
		request.PostForm = map[string][]string{
			"title":   {"Title"},
			"content": {"Content"},
			"expires": {"1"},
		}

		response := httptest.NewRecorder()

		// Adding ExpectPrepare to DB Expectations
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", "2024-01-23T10:23:42Z", "2024-01-24T10:23:42Z")
		prep.ExpectExec().WithArgs(
			"Title",
			"Content",
			"1",
		).WillReturnResult(sqlmock.NewResult(0, 1))

		server.Handler.ServeHTTP(response, request)

		// It now redirects to another page. I should continue reading the book for more info.
		assertStatus(t, response, http.StatusSeeOther)
	})
	// We are now showing the form which allows the user to enter data to be POST-ed
	t.Run("checking create snippet OK Case - GET instead of POST", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodGet, "snippet/create")
		response := httptest.NewRecorder()

		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusOK)
	})
	t.Run("checking create snippet NOK Case - malformed URL", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodGet, "snippet/create/?id=1")
		response := httptest.NewRecorder()

		// Adding ExpectPrepare to DB Expectations
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", "2024-01-23T10:23:42Z", "2024-01-24T10:23:42Z")
		prep.ExpectQuery().WithArgs(1).WillReturnRows(rows)

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
