package server

import (
	"fmt"
	"net/http"
)

var StaticFolder = "./ui/static"

// NOTE:
// Set homeTemplate and showSnippetTemplate to nil
// if you want to use the default values
func CreateServer(app *Application,
	homeTemplate, showSnippetTemplate []string) (*http.Server, error) {

	// Initializing the Template Files
	if homeTemplate != nil {
		homePageTemplateFiles = homeTemplate
	}
	if showSnippetTemplate != nil {
		showSnippetTemplateFiles = showSnippetTemplate
	}

	fileServer := createFileServer()
	routes, err := createRoutes(fileServer, app)
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

func createRoutes(fileServer http.Handler, app *Application) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	if mux != nil {
		mux.HandleFunc("/", app.home)
		mux.HandleFunc("/snippet", app.showSnippet)
		mux.HandleFunc("/snippet/create", app.createSnippet)
		mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	} else {
		return nil, fmt.Errorf("cannot create Handler")
	}
	return mux, nil
}

func createFileServer() http.Handler {
	return http.FileServer(http.Dir(StaticFolder))
}
