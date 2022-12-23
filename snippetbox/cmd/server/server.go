package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var TemplateFiles = []string{
	"./ui/html/home.page.tmpl",
	"./ui/html/base.layout.tmpl",
	"./ui/html/footer.partial.tmpl",
}

var StaticFolder = "./ui/static"

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func CreateServer(addr *string, errorLog *log.Logger, infoLog *log.Logger) (*http.Server, error) {
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}
	fileServer := createFileServer()
	routes, err := createRoutes(fileServer, app)
	if err != nil {
		msg := "creating routes failed"
		errorLog.Println(msg)
		return nil, fmt.Errorf(msg)
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  routes,
	}
	return srv, nil
}

func createRoutes(fileServer http.Handler, app *application) (*http.ServeMux, error) {
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

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ts, err := template.ParseFiles(TemplateFiles...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Template file not found", http.StatusNotFound)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error: Cannot execute Template", http.StatusInternalServerError)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
