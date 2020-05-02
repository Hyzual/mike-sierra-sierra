package main

import (
	"log"
	"net/http"

	"github.com/hyzual/mike-sierra-sierra/server"
)

func main() {
	handler := http.HandlerFunc(server.MusicServer)
	err := http.ListenAndServe(":80", handler)
	if err != nil {
		log.Fatalf("could not listen on port 80 %v", err)
	}
}
