package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hyzual/mike-sierra-sierra/server"
)

func TestGetRoot(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	server.MusicServer(response, request)

	got := response.Code
	want := 200
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
