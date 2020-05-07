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
)

// MusicServer serves HTTP requests.
// It serves the HTML pages, the REST routes and the media files as well.
type MusicServer struct {
	assetsIncluder AssetsIncluder
	templateLoader TemplateLoader
	http.Handler
}

// New creates a new MusicServer
func New(assetsIncluder AssetsIncluder, templateLoader TemplateLoader, musicLoader PathJoiner, loginHandler *user.LoginHandler) *MusicServer {
	s := new(MusicServer)
	s.assetsIncluder = assetsIncluder
	s.templateLoader = templateLoader

	router := mux.NewRouter()
	router.HandleFunc("/", s.rootHandler)
	router.HandleFunc("/home", s.homeHandler)
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

func (s *MusicServer) homeHandler(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := s.templateLoader.Load("app.html")
	if err != nil {
		http.Error(writer, fmt.Sprintf("could not load template %s", err), http.StatusInternalServerError)
	}
	err = tmpl.Execute(writer, nil)
	if err != nil {
		http.Error(writer, fmt.Sprintf("could not execute template %s", err), http.StatusInternalServerError)
	}
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
	//TODO: I don't know if it's worth it to clean that path. I never get the paths with dots
	cleanedPath := path.Clean(request.URL.Path)
	http.ServeFile(writer, request, m.musicLoader.Join(cleanedPath))
}
