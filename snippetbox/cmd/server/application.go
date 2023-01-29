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
	Addr            *string
	InfoLog         *log.Logger
	ErrorLog        *log.Logger
	SnippetModel    *mysql.SnippetModel
	SnippetContents *models.SnippetContents
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

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Println("----- Inside Home() ---- ")
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	s, err := app.SnippetModel.Latest()
	if err != nil {
		app.ErrorLog.Printf("\n\t-----Home(): Error Found: %s -----", err)
		app.serverError(w, err)
		return
	}

	for _, snippet := range s {
		fmt.Fprintf(w, "%v\n", snippet)
	}

	ts, err := template.ParseFiles(homePageTemplateFiles...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.SnippetModel.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, s)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) createSnippet(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Println("----- Inside CreateSnippet() ---- ")
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	if app.SnippetContents == nil {
		app.ErrorLog.Println("Sql Record is not defined")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if app.SnippetModel == nil {
		app.ErrorLog.Println("Snippets is not defined")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := app.SnippetModel.Insert(
		app.SnippetContents.Title,
		app.SnippetContents.Content,
		app.SnippetContents.Expires)
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
