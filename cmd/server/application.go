package server

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
	"snippetbox/pkg/forms"
	"snippetbox/pkg/models"
	"snippetbox/pkg/models/mysql"
	"strconv"

	"github.com/justinas/alice"

	"github.com/bmizerany/pat"

	"strings"
	"unicode/utf8"

	"github.com/golangcollege/sessions"
)

var StaticFolder = "./ui/static"

type Application struct {
	Port          *string
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Snippets      *mysql.SnippetDatabase
	TemplateCache map[string]*template.Template
	Session       *sessions.Session
	TLSConfig     *tls.Config
}

var homePageTemplateFiles = []string{
	"./ui/html/home.page.tmpl",
	"./ui/html/base.layout.tmpl",
	"./ui/html/footer.partial.tmpl",
}

var showSnippetTemplateFiles = []string{
	"./ui/html/show.page.tmpl",
	"./ui/html/base.layout.tmpl",
	"./ui/html/footer.partial.tmpl",
}

func CreateServer(app *Application) (*http.Server, error) {
	routes, err := app.createRoutes()
	if err != nil {
		msg := "creating routes failed"
		app.ErrorLog.Println(msg)
		return nil, fmt.Errorf(msg)
	}

	srv := &http.Server{
		Addr:      *app.Port,
		ErrorLog:  app.ErrorLog,
		Handler:   routes,
		TLSConfig: app.TLSConfig,
	}
	return srv, nil
}

func (app *Application) createRoutes() (http.Handler, error) {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.Session.Enable)
	mux := pat.New()

	// Update these routes to use the new dynamic middleware chain followed
	// by the appropriate handler function.
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))
	mux.Get("/{.*}", dynamicMiddleware.ThenFunc(app.notFound))

	// Leave the static files route unchanged.
	fileServer := http.FileServer(http.Dir(StaticFolder))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux), nil
}

func (app *Application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Printf("\n\t home() called")
	s, err := app.Snippets.Latest()
	if err != nil {
		app.ErrorLog.Printf("\n\tError: %s", err)
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippet:  nil,
		Snippets: s,
	})
}

func (app *Application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.ErrorLog.Printf("\n\tError: %s", err)
		app.notFound(w, r)
		return
	}

	snippet, err := app.Snippets.Get(id)
	switch {
	case err == models.ErrNoRecord:
		app.ErrorLog.Printf("\n\tError: %s", err)
		app.notFound(w, r)
		return
	case err != nil:
		app.ErrorLog.Printf("\n\tError: %s", err)
		app.serverError(w, err)
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: snippet,
	})
}

func (app *Application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	id, err := app.Snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.Session.Put(r, "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func validateSnippets(title, content, expires string) map[string]string {
	errors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This field is too long (maximum is 100 characters)"
	}

	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field cannot be blank"
	}

	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "This field cannot be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "This field is invalid"
	}
	return errors
}

func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// NOTE: r *http.Request is ignored by this function.
// It's only used so that I can attach this to the routes
func (app *Application) notFound(w http.ResponseWriter, r *http.Request) {
	app.ErrorLog.Printf("%s not found", r.URL.Path)
	app.clientError(w, http.StatusNotFound)
}
