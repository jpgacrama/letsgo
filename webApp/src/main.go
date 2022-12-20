package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from Snippetbox"))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a new snippet..."))
}

func createRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	if mux != nil {
		mux.HandleFunc("/", home)
		mux.HandleFunc("/snippet", showSnippet)
		mux.HandleFunc("/snippet/create", createSnippet)
	} else {
		log.Fatal("Cannot create Server MUX")
	}
	return mux
}

func main() {
	log.Println("Starting server on :4000")
	mux := createRoutes()
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
