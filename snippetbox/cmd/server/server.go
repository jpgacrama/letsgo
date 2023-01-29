package server

import (
	"fmt"
	"net/http"
)

var StaticFolder = "./ui/static"

func CreateServer(app *Application, tmplFiles ...string) (*http.Server, error) {
	if len(tmplFiles) > 0 {
		homePageTemplateFiles = tmplFiles
	}
	fileServer := createFileServer()
	routes, err := createRoutes(fileServer, app)
	if err != nil {
		msg := "creating routes failed"
		app.ErrorLog.Println(msg)
		return nil, fmt.Errorf(msg)
	}

	srv := &http.Server{
		Addr:     *app.Addr,
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
