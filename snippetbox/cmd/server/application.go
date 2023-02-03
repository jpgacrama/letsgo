package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
	"snippetbox/pkg/models"
	"snippetbox/pkg/models/mysql"
	"strconv"
)

type Application struct {
	Port          *string
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	SnippetDB     *mysql.SnippetDatabase
	Snippet       *models.Snippet
	TemplateCache map[string]*template.Template
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

var StaticFolder = "./ui/static"

func CreateServer(app *Application) (*http.Server, error) {
	fileServer := createFileServer()
	routes, err := app.createRoutes(fileServer)
	if err != nil {
		msg := "creating routes failed"
		app.ErrorLog.Println(msg)
		return nil, fmt.Errorf(msg)
	}

	srv := &http.Server{
		Addr:     *app.Port,
		ErrorLog: app.ErrorLog,
		Handler:  routes,
	}
	return srv, nil
}

func createFileServer() http.Handler {
	return http.FileServer(http.Dir(StaticFolder))
}

func (app *Application) createRoutes(fileServer http.Handler) (http.Handler, error) {
	mux := http.NewServeMux()
	if mux != nil {
		mux.HandleFunc("/", app.home)
		mux.HandleFunc("/snippet", app.showSnippet)
		mux.HandleFunc("/snippet/create", app.createSnippet)
		mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	} else {
		return nil, fmt.Errorf("cannot create Handler")
	}

	return app.recoverPanic(app.logRequest(secureHeaders(mux))), nil
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.ErrorLog.Printf("\n\t----- home() error. Path is: %s -----", r.URL.Path)
		app.notFound(w)
		return
	}

	s, err := app.SnippetDB.Latest()
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
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.ErrorLog.Printf("\n\tError: %s", err)
		app.notFound(w)
		return
	}

	snippetContents, err := app.SnippetDB.Get(id)
	switch {
	case err == models.ErrNoRecord:
		app.ErrorLog.Printf("\n\tError: %s", err)
		app.notFound(w)
		return
	case err != nil:
		app.ErrorLog.Printf("\n\tError: %s", err)
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: snippetContents,
	})
}

func (app *Application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	switch {
	case app.Snippet == nil:
		app.ErrorLog.Println("Sql Record is not defined")
		app.clientError(w, http.StatusBadRequest)
		return
	case app.SnippetDB == nil:
		app.ErrorLog.Println("Snippet Database is not defined")
		app.clientError(w, http.StatusBadRequest)
		return

	}

	id, err := app.SnippetDB.Insert(
		app.Snippet.Title,
		app.Snippet.Content,
		app.Snippet.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
