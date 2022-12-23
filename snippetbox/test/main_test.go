package main_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"snippetbox/cmd/server"
	"testing"
)

var templateFiles = []string{
	"../ui/html/home.page.tmpl",
	"../ui/html/base.layout.tmpl",
	"../ui/html/footer.partial.tmpl",
}

var staticFolder = "../ui/static"

func TestHomePage(t *testing.T) {
	server.TemplateFiles = templateFiles
	t.Run("checking home page OK Case", func(t *testing.T) {
		server, err := server.CreateServer()
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest("")
		response := httptest.NewRecorder()
		server.GetMux().ServeHTTP(response, request)
		assertStatus(t, response, http.StatusOK)
	})
	t.Run("checking home page NOK Case", func(t *testing.T) {
		server, err := server.CreateServer()
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest("123")
		response := httptest.NewRecorder()
		server.GetMux().ServeHTTP(response, request)
		assertStatus(t, response, http.StatusNotFound)
	})
}

func TestStaticPage(t *testing.T) {
	server.TemplateFiles = templateFiles
	server.StaticFolder = staticFolder
	t.Run("checking static page OK Case", func(t *testing.T) {
		server, err := server.CreateServer()
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest("static/")
		response := httptest.NewRecorder()
		server.GetMux().ServeHTTP(response, request)
		assertStatus(t, response, http.StatusOK)
	})
	t.Run("checking static page NOK Case", func(t *testing.T) {
		server, err := server.CreateServer()
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest("static/123")
		response := httptest.NewRecorder()
		server.GetMux().ServeHTTP(response, request)
		assertStatus(t, response, http.StatusNotFound)
	})
}

func newRequest(str string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", str), nil)
	return req
}

func assertStatus(t testing.TB, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	got := response.Code
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
