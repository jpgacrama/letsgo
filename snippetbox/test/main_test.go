package main_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"snippetbox/cmd/handlers"
	"testing"
)

func TestHandler(t *testing.T) {
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
}

func newRequest(str string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", str), nil)
	return req

}

func assertStatus(t testing.TB, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	got := response.Code
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
