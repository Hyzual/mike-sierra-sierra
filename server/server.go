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
	"html/template"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/hyzual/mike-sierra-sierra/server/user"
)

// MusicServer serves HTTP requests.
// It serves the HTML pages, the REST routes and the media files as well.
type MusicServer struct {
	pathJoiner PathJoiner
	http.Handler
}

// New creates a new MusicServer
func New(pathJoiner PathJoiner, loginHandler *user.LoginHandler) *MusicServer {
	s := new(MusicServer)
	s.pathJoiner = pathJoiner

	router := mux.NewRouter()
	router.HandleFunc("/", s.rootHandler)
	router.HandleFunc("/home", s.homeHandler)
	router.HandleFunc("/login", s.getLoginHandler).Methods(http.MethodGet)
	router.Handle("/login", loginHandler).Methods(http.MethodPost)
	router.PathPrefix("/assets/").HandlerFunc(s.assetsHandler)

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
	fmt.Fprint(writer, "Hello world")
}

func (s *MusicServer) getLoginHandler(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles(s.pathJoiner.Join("./templates/login.html"))

	if err != nil {
		http.Error(writer, fmt.Sprintf("problem loading template %s", err.Error()), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(writer, nil)
	if err != nil {
		http.Error(writer, fmt.Sprintf("problem executing template %s", err.Error()), http.StatusInternalServerError)
	}
}

func (s *MusicServer) assetsHandler(writer http.ResponseWriter, request *http.Request) {
	cleanedPath := path.Clean(request.URL.Path)
	if cleanedPath == "/assets" {
		http.NotFound(writer, request)
		return
	}

	http.ServeFile(writer, request, s.pathJoiner.Join(cleanedPath))
}
