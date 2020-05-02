package server

import (
	"fmt"
	"net/http"
)

// MusicServer serves HTTP requests.
// It serves the HTML pages as well as the media files.
func MusicServer(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Hello world")
}
