package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

	// Execute the template set, passing in any dynamic data.
	err := ts.Execute(w, td)
	if err != nil {
		app.ErrorLog.Printf("\n\t--- render() error; %s ---", err)
		app.serverError(w, err)
	}
}
