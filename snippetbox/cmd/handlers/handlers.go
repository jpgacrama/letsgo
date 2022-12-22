package handlers

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

type Server struct {
	handler http.Handler
}

func CreateServer() (*Server, error) {
	p := new(Server)
	router := http.NewServeMux()
	if router != nil {
		router.HandleFunc("/", home)
		router.HandleFunc("/snippet", showSnippet)
		router.HandleFunc("/snippet/create", createSnippet)
	} else {
		return nil, fmt.Errorf("cannot create Handler")
	}
	p.handler = router
	return p, nil
}

func (s *Server) GetHandler() http.Handler {
	return s.handler
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
