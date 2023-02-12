package snippetbox_test

import (
	"flag"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golangcollege/sessions"
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
	session := sessions.New([]byte(*createSession()))
	session.Lifetime = 12 * time.Hour

	app := &server.Application{
		Port:     &port,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Session:  session,
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
	t.Run("checking show snippet NOK Case - malformed URL", func(t *testing.T) {
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
	t.Run("checking show snippet NOK Case - no ID found", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodGet, "snippet/10")
		response := httptest.NewRecorder()

		// Adding ExpectPrepare to DB Expectations
		rows := sqlmock.NewRows([]string{})
		prep.ExpectQuery().WithArgs(10).WillReturnRows(rows)
		prep.ExpectQuery().WithArgs().WillReturnRows(rows)

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
	t.Run("checking create snippet NOK Case - Title is too long", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodPost, "snippet/create")

		tooLongTitle := `Lorem Ipsum is simply dummy text of the printing and typesetting industry.
		Lorem Ipsum has been the industry's standard dummy text ever since the 1500s,
		when an unknown printer took a galley of type and scrambled it to make a type specimen book.
		It has survived not only five centuries, but also the leap into electronic typesetting,
		remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset
		sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like
		Aldus PageMaker including versions of Lorem Ipsum.

		Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of
		classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a
		Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin
		words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in
		classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32
		and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero,
		written in 45 BC. This book is a treatise on the theory of ethics, very popular during the
		Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a
		line in section 1.10.32.

		The standard chunk of Lorem Ipsum used since the 1500s is reproduced below for those interested.
		Sections 1.10.32 and 1.10.33 from "de Finibus Bonorum et Malorum" by Cicero are also reproduced
		in their exact original form, accompanied by English versions from the 1914 translation by H. Rackham.

		There are many variations of passages of Lorem Ipsum available, but the majority have suffered
		alteration in some form, by injected humour, or randomised words which don't look even slightly
		believable. If you are going to use a passage of Lorem Ipsum, you need to be sure there isn't
		anything embarrassing hidden in the middle of text. All the Lorem Ipsum generators on the Internet
		tend to repeat predefined chunks as necessary, making this the first true generator on the Internet.
		It uses a dictionary of over 200 Latin words, combined with a handful of model sentence structures,
		to generate Lorem Ipsum which looks reasonable. The generated Lorem Ipsum is therefore always free
		from repetition, injected humour, or non-characteristic words etc.`

		request.PostForm = map[string][]string{
			"title":   {tooLongTitle},
			"content": {"Content"},
			"expires": {"1"},
		}
		response := httptest.NewRecorder()

		// Adding ExpectPrepare to DB Expectations
		prep.ExpectExec().WithArgs(
			tooLongTitle,
			"Content",
			"1",
		).WillReturnResult(sqlmock.NewResult(0, 0))

		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusOK) // Error message is displayed on screen instead
	})
	t.Run("checking create snippet NOK Case - Title is blank", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodPost, "snippet/create")
		blankTitle := ""
		request.PostForm = map[string][]string{
			"title":   {blankTitle},
			"content": {"Content"},
			"expires": {"1"},
		}
		response := httptest.NewRecorder()

		// Adding ExpectPrepare to DB Expectations
		prep.ExpectExec().WithArgs(
			blankTitle,
			"Content",
			"1",
		).WillReturnResult(sqlmock.NewResult(0, 0))

		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusOK) // Error message is displayed on screen instead
	})
	t.Run("checking create snippet NOK Case - Content is blank", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodPost, "snippet/create")
		blankContent := ""
		request.PostForm = map[string][]string{
			"title":   {"Title"},
			"content": {blankContent},
			"expires": {"1"},
		}
		response := httptest.NewRecorder()

		// Adding ExpectPrepare to DB Expectations
		prep.ExpectExec().WithArgs(
			"Title",
			blankContent,
			"1",
		).WillReturnResult(sqlmock.NewResult(0, 0))

		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusOK) // Error message is displayed on screen instead
	})

	t.Run("checking create snippet NOK Case - Expires is blank", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodPost, "snippet/create")
		blankExpires := ""
		request.PostForm = map[string][]string{
			"title":   {"Title"},
			"content": {"Content"},
			"expires": {blankExpires},
		}
		response := httptest.NewRecorder()

		// Adding ExpectPrepare to DB Expectations
		prep.ExpectExec().WithArgs(
			"Title",
			"Content",
			blankExpires,
		).WillReturnResult(sqlmock.NewResult(0, 0))

		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusOK) // Error message is displayed on screen instead
	})
	// NOK Case where Expires is not any of these values: 365, 7 , or 1
	t.Run("checking create snippet NOK Case - Expires value is invalid", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest(http.MethodPost, "snippet/create")
		wrongExpiresValue := "25"
		request.PostForm = map[string][]string{
			"title":   {"Title"},
			"content": {"Content"},
			"expires": {wrongExpiresValue},
		}
		response := httptest.NewRecorder()

		// Adding ExpectPrepare to DB Expectations
		prep.ExpectExec().WithArgs(
			"Title",
			"Content",
			wrongExpiresValue,
		).WillReturnResult(sqlmock.NewResult(0, 0))

		server.Handler.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusOK) // Error message is displayed on screen instead
	})
}

func TestCatchAll(t *testing.T) {
	db, mock := NewMock()
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
	t.Run("checking catch-all", func(t *testing.T) {
		server, err := server.CreateServer(app)
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}

		// Adding ExpectPrepare to DB Expectations
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", time.Now(), "2024-01-24T10:23:42Z")
		prep.ExpectQuery().WillReturnRows(rows)

		request := newRequest(http.MethodGet, "jonas")
		wrongExpiresValue := "25"
		request.PostForm = map[string][]string{
			"title":   {""},
			"content": {""},
			"expires": {wrongExpiresValue},
		}
		response := httptest.NewRecorder()
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

func createSession() *string {
	// TODO: Change the secret string to your choice
	secret := new(string)
	if !flag.Parsed() {
		secret = flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
		flag.Parse()
	}
	return secret
}
