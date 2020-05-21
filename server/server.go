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
Package server implements the music server.
It handles all HTTP routing, serves HTML pages, REST routes and media files.
*/
package server

import (
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/swithek/sessionup"
)

// MusicServer serves HTTP requests.
// It serves the HTML pages, the REST routes and the media files as well.
type MusicServer struct {
	sessionManager *sessionup.Manager
	assetsLoader   PathJoiner
	assetsResolver AssetsResolver
	http.Handler
}

// New creates a new MusicServer
func New(
	router *mux.Router,
	sessionManager *sessionup.Manager,
	assetsLoader PathJoiner,
	templateExecutor TemplateExecutor,
	musicLoader PathJoiner,
	assetsResolver AssetsResolver,
) *MusicServer {
	s := new(MusicServer)
	s.sessionManager = sessionManager
	s.assetsLoader = assetsLoader

	homeHandler := sessionManager.Auth(
		WrapErrors(&homeHandler{templateExecutor, assetsResolver}),
	)
	musicHandler := &musicHandler{musicLoader}

	router.HandleFunc("/", s.rootHandler)
	router.Handle("/home", homeHandler)
	router.PathPrefix("/assets/").HandlerFunc(s.assetsHandler)
	router.PathPrefix("/music/").Handler(http.StripPrefix("/music/", musicHandler))

	s.Handler = router
	return s
}

func (s *MusicServer) rootHandler(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/home", http.StatusFound)
}

func (s *MusicServer) assetsHandler(writer http.ResponseWriter, request *http.Request) {
	cleanedPath := path.Clean(request.URL.Path)
	if cleanedPath == "/assets" {
		http.NotFound(writer, request)
		return
	}

	http.ServeFile(writer, request, s.assetsLoader.Join(cleanedPath))
}

type musicHandler struct {
	pathJoiner PathJoiner
}

func (m *musicHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, m.pathJoiner.Join(request.URL.Path))
}

// HandleUnauthorized redirects to /login when users are not authenticated.
// It is used by sessionup's Auth middleware.
func HandleUnauthorized(_ error) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "/login", http.StatusFound)
	})
}
