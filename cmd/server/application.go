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
	"github.com/golangcollege/sessions"
	"time"
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
	Users         *mysql.UserModel
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
	routes := app.createRoutes()
	srv := &http.Server{
		Addr:         *app.Port,
		ErrorLog:     app.ErrorLog,
		Handler:      routes,
		TLSConfig:    app.TLSConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return srv, nil
}

func (app *Application) createRoutes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.Session.Enable)
	mux := pat.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))
	mux.Get("/snippet?id=1", dynamicMiddleware.ThenFunc(app.showSnippet))
	mux.Get("/{.*}", dynamicMiddleware.ThenFunc(app.notFound))

	// New routes for signing-in and out, and authentication
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir(StaticFolder))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
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
		app.badRequest(w, r)
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
	if err := r.ParseForm(); err != nil {
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

func (app *Application) badRequest(w http.ResponseWriter, r *http.Request) {
	app.ErrorLog.Printf("Bad Request: %s", r.URL.Path)
	app.clientError(w, http.StatusBadRequest)
}

func (app *Application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *Application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	err = app.Users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.Session.Put(r, "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *Application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *Application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, err := app.Users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// User is now logged in at this point
	app.Session.Put(r, "userID", id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *Application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.Session.Remove(r, "userID")
	app.Session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
