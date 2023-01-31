package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func CreateLoggers() (*log.Logger, *log.Logger) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return infoLog, errorLog
}

// Retrieve the appropriate template set from the cache based on the page name
// (like 'home.page.tmpl'). If no entry exists in the cache with the
// provided name, call the serverError helper method that we made earlier.
func (app *Application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.TemplateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.ErrorLog.Printf("\n\t--- render() error: %s ---", err)
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

// This takes a pointer to a templateData struct, adds the current year
// to the CurrentYear field, and then returns the pointer.
// Again, we're not using the *http.Request parameter at the
// moment, but we will do later in the book.
func (app *Application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	return td
}
