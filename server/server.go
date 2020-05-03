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
	"path/filepath"
)

// MusicServer serves HTTP requests.
// It serves the HTML pages, the REST routes and the media files as well.
type MusicServer struct {
	rootDirectory string //Project root directory from which to load templates and assets as a relative path
	http.Handler
}

// New creates a new MusicServer
func New(rootDirectory string) *MusicServer {
	s := new(MusicServer)
	s.rootDirectory = rootDirectory

	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(s.rootHandler))
	router.Handle("/home", http.HandlerFunc(s.homeHandler))
	router.Handle("/login", http.HandlerFunc(s.loginHandler))
	router.Handle("/assets/", http.HandlerFunc(s.assetsHandler))

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

func (s *MusicServer) loginHandler(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join(s.rootDirectory, "./templates/login.html"))

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
		http.Error(writer, "Forbidden", http.StatusForbidden)
		return
	}

	joinedPath := path.Join(s.rootDirectory, cleanedPath)

	http.ServeFile(writer, request, joinedPath)
}
