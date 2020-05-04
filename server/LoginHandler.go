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

package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

// LoginHandler handles GET and POST to /login route
type LoginHandler struct {
	pathJoiner PathJoiner
	userStore  UserStore
}

// NewLoginHandler create a new LoginHandler
func NewLoginHandler(pathJoiner PathJoiner, userStore UserStore) *LoginHandler {
	return &LoginHandler{pathJoiner, userStore}
}

func (l *LoginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		l.loginGetHandler(writer, request)
		return
	}

	if request.Method == http.MethodPost {
		l.loginPostHandler(writer, request)
		return
	}

	http.NotFound(writer, request)
}

func (l *LoginHandler) loginGetHandler(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles(l.pathJoiner.Join("./templates/login.html"))

	if err != nil {
		http.Error(writer, fmt.Sprintf("problem loading template %s", err.Error()), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(writer, nil)
	if err != nil {
		http.Error(writer, fmt.Sprintf("problem executing template %s", err.Error()), http.StatusInternalServerError)
	}
}

func (l *LoginHandler) loginPostHandler(writer http.ResponseWriter, request *http.Request) {
	decoder := schema.NewDecoder()
	err := request.ParseForm()
	if err != nil {
		log.Println(errors.Wrap(err, "Could not parse the login form"))
		badRequest(writer, request)
		return
	}

	loginForm := new(LoginFormRepresentation)
	decoder.ZeroEmpty(true)
	err = decoder.Decode(loginForm, request.PostForm)
	if err != nil {
		log.Println(errors.Wrap(err, "Could not decode the login form into its representation"))
		http.Error(writer, "Bad Request", http.StatusBadRequest)
		return
	}

	err = l.userStore.VerifyCredentialsMatch(*loginForm)
	if err != nil {
		http.Error(writer, "Forbidden", http.StatusForbidden)
		return
	}
	http.Redirect(writer, request, "/home", http.StatusFound)
}

// badRequest replies to the request with an HTTP 400 Bad request error.
func badRequest(writer http.ResponseWriter, request *http.Request) {
	http.Error(writer, "Bad Request", http.StatusBadRequest)
}

// LoginFormRepresentation represents the login form that is expected from the HTML
type LoginFormRepresentation struct {
	Email    string `schema:"email,required"`
	Password string `schema:"password,required"`
}
