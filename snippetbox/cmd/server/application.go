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
	Addr     *string
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	DB       *mysql.SnippetDatabase
	Snippet  *models.Snippet
}

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Println("----- Inside Home() ---- ")
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	ts, err := template.ParseFiles(templateFiles...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) ShowSnippet(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Println("----- Inside ShowSnippet() ---- ")
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	result, err := app.DB.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	fmt.Fprintf(w, "%v", result)
}

func (app *Application) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Println("----- Inside CreateSnippet() ---- ")
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	if app.Snippet == nil {
		app.ErrorLog.Println("Sql Record is not defined")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if app.DB == nil {
		app.ErrorLog.Println("Snippets is not defined")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := app.DB.Insert(
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
