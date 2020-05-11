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
	"fmt"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/hyzual/mike-sierra-sierra/server/user"
	"github.com/swithek/sessionup"
)

// MusicServer serves HTTP requests.
// It serves the HTML pages, the REST routes and the media files as well.
type MusicServer struct {
	sessionManager *sessionup.Manager
	assetsIncluder AssetsIncluder
	templateLoader TemplateLoader
	http.Handler
}

// New creates a new MusicServer
func New(sessionManager *sessionup.Manager, assetsIncluder AssetsIncluder, templateLoader TemplateLoader, musicLoader PathJoiner, loginHandler *user.LoginHandler) *MusicServer {
	s := new(MusicServer)
	s.sessionManager = sessionManager
	s.assetsIncluder = assetsIncluder
	s.templateLoader = templateLoader

	homeHandler := NewHomeHandler(templateLoader)

	router := mux.NewRouter()
	router.HandleFunc("/", s.rootHandler)
	router.Handle("/home", sessionManager.Auth(homeHandler))
	router.HandleFunc("/login", s.getLoginHandler).Methods(http.MethodGet)
	router.Handle("/login", loginHandler).Methods(http.MethodPost)
	router.PathPrefix("/assets/").HandlerFunc(s.assetsHandler)

	musicHandler := &musicHandler{musicLoader}
	router.PathPrefix("/music/").Handler(http.StripPrefix("/music/", musicHandler))

	s.Handler = router
	return s
}

func (s *MusicServer) rootHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/" {
		http.Redirect(writer, request, "/home", http.StatusFound)
		return
	}
	http.NotFound(writer, request)
}

func (s *MusicServer) getLoginHandler(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := s.templateLoader.Load("login.html")
	if err != nil {
		http.Error(writer, fmt.Sprintf("could not load template %s", err), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(writer, nil)
	if err != nil {
		http.Error(writer, fmt.Sprintf("could not execute template %s", err), http.StatusInternalServerError)
	}
}

func (s *MusicServer) assetsHandler(writer http.ResponseWriter, request *http.Request) {
	cleanedPath := path.Clean(request.URL.Path)
	if cleanedPath == "/assets" {
		http.NotFound(writer, request)
		return
	}

	http.ServeFile(writer, request, s.assetsIncluder.Join(cleanedPath))
}

type musicHandler struct {
	musicLoader PathJoiner
}

func (m *musicHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, m.musicLoader.Join(request.URL.EscapedPath()))
}

// HomeHandler handles GET /home. It renders the app template and inits the app
type HomeHandler struct {
	templateLoader TemplateLoader
}

// NewHomeHandler creates a new HomeHandler
func NewHomeHandler(templateLoader TemplateLoader) http.Handler {
	return &HomeHandler{templateLoader}
}

func (h *HomeHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := h.templateLoader.Load("app.html")
	if err != nil {
		http.Error(writer, fmt.Sprintf("could not load template %s", err), http.StatusInternalServerError)
	}
	err = tmpl.Execute(writer, nil)
	if err != nil {
		http.Error(writer, fmt.Sprintf("could not execute template %s", err), http.StatusInternalServerError)
	}
}

// HandleUnauthorized redirects to /login when users are not authenticated
// It is used by sessionup's Auth middleware
func HandleUnauthorized(_ error) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "/login", http.StatusFound)
	})
}
