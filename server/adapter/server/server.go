/*
 *   Copyright (C) 2020-2021  Joris MASSON
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU Affero General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU Affero General Public License for more details.
 *
 *   You should have received a copy of the GNU Affero General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

/*
Package server implements the music server.
It handles all HTTP routing, serves HTML pages, REST routes and media files.
*/
package server

import (
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/hyzual/mike-sierra-sierra/server/adapter"
)

// MusicServer serves HTTP requests.
// It serves the HTML pages, the REST routes and the media files as well.
type MusicServer struct {
	http.Handler
}

// New creates a new MusicServer
func New(
	router *mux.Router,
	assetsLoader adapter.PathJoiner,
	musicLoader adapter.PathJoiner,
) *MusicServer {
	s := new(MusicServer)

	musicHandler := &musicHandler{musicLoader}
	assetsHandler := &assetsHandler{assetsLoader}

	router.HandleFunc("/", s.rootHandler)
	router.PathPrefix("/assets/").Handler(assetsHandler)
	router.PathPrefix("/music/").Handler(http.StripPrefix("/music/", musicHandler))

	s.Handler = router
	return s
}

func (s *MusicServer) rootHandler(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/app", http.StatusFound)
}

type assetsHandler struct {
	assetsLoader adapter.PathJoiner
}

func (a *assetsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	cleanedPath := path.Clean(request.URL.Path)
	if cleanedPath == "/assets" {
		http.NotFound(writer, request)
		return
	}

	http.ServeFile(writer, request, a.assetsLoader.Join(cleanedPath))
}

type musicHandler struct {
	pathJoiner adapter.PathJoiner
}

func (m *musicHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, m.pathJoiner.Join(request.URL.Path))
}

// HandleUnauthorized redirects to /sign-in when users are not authenticated.
// It is used by sessionup's Auth middleware.
func HandleUnauthorized(_ error) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "/sign-in", http.StatusFound)
	})
}
