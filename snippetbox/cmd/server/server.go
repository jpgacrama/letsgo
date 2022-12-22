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

type Server struct {
	mux        *http.ServeMux
	fileServer http.Handler
}

func CreateServer() (*Server, error) {
	p := new(Server)
	fileServer := createFileServer()
	routes, err := createRoutes(fileServer)
	if err != nil {
		log.Fatalln("creating routes failed")
	}
	p.mux = routes
	p.fileServer = fileServer
	return p, nil
}

func createRoutes(fileServer http.Handler) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	if mux != nil {
		mux.HandleFunc("/", home)
		mux.HandleFunc("/snippet", showSnippet)
		mux.HandleFunc("/snippet/create", createSnippet)
	} else {
		return nil, fmt.Errorf("cannot create Handler")
	}
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux, nil
}

func createFileServer() http.Handler {
	return http.FileServer(http.Dir(StaticFolder))
}

func (s *Server) GetMux() *http.ServeMux {
	return s.mux
}

func home(w http.ResponseWriter, r *http.Request) {
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

func showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
