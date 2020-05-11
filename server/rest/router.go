/*
 *   Copyright (c) 2020 Joris MASSON

 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.

 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.

 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

/*
Package rest implements the REST API.
*/
package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swithek/sessionup"
)

// Register registers a gorilla/mux Subrouter for the REST API on the given router
func Register(router *mux.Router, sessionManager *sessionup.Manager) {
	songHandler := &songHandler{}

	apiRouter := router.PathPrefix("/api/").Subrouter()
	apiRouter.Handle("/songs/{songId}", sessionManager.Auth(songHandler))
}

// Song represents a song's details
type Song struct {
	Title string `json:"title"`
}

const contentTypeHeader = "Content-Type"
const jsonMediaType = "application/json; charset=utf-8"

type songHandler struct {
}

func (s *songHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	resp := &Song{Title: "Hello World"}
	json.NewEncoder(writer).Encode(resp)

	writer.Header().Set("Content-Type", jsonMediaType)
	writer.WriteHeader(http.StatusOK)
}
