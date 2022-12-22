package main_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"snippetbox/cmd/handlers"
	"testing"
)

func TestHomePage(t *testing.T) {
	t.Run("checking home page OK Case", func(t *testing.T) {
		handlers.TemplateFile = "../ui/html/home.page.tmpl"
		server, err := handlers.CreateServer()
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest("")
		response := httptest.NewRecorder()
		server.GetHandler().ServeHTTP(response, request)
		assertStatus(t, response, http.StatusOK)
	})
	t.Run("checking home page NOK Case", func(t *testing.T) {
		handlers.TemplateFile = "../ui/html/home.page.tmpl"
		server, err := handlers.CreateServer()
		if err != nil {
			log.Fatalf("problem creating server %v", err)
		}
		request := newRequest("123")
		response := httptest.NewRecorder()
		server.GetHandler().ServeHTTP(response, request)
		assertStatus(t, response, http.StatusNotFound)
	})
}

// func TestShowSnippet(t *testing.T) {
// 	t.Run("checking show snippet OK Case", func(t *testing.T) {
// 		server, err := handlers.CreateServer()
// 		if err != nil {
// 			log.Fatalf("problem creating server %v", err)
// 		}
// 		request := newRequest("/snippet")
// 		response := httptest.NewRecorder()
// 		server.GetHandler().ServeHTTP(response, request)
// 		assertStatus(t, response, http.StatusOK)
// 	})
// 	t.Run("checking show snippet NOK Case", func(t *testing.T) {
// 		server, err := handlers.CreateServer()
// 		if err != nil {
// 			log.Fatalf("problem creating server %v", err)
// 		}
// 		request := newRequest("123")
// 		response := httptest.NewRecorder()
// 		server.GetHandler().ServeHTTP(response, request)
// 		assertStatus(t, response, http.StatusNotFound)
// 	})
// }

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
