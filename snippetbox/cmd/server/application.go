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
	Addr      *string
	InfoLog   *log.Logger
	ErrorLog  *log.Logger
	Snippets  *mysql.SnippetModel
	SqlRecord *models.Record
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
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *Application) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Println("----- Inside CreateSnippet() ---- ")
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	if app.SqlRecord == nil {
		app.ErrorLog.Println("Sql Record is not defined")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if app.Snippets == nil {
		app.ErrorLog.Println("Snippets is not defined")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := app.Snippets.Insert(
		app.SqlRecord.Title,
		app.SqlRecord.Content,
		app.SqlRecord.Expires)
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
