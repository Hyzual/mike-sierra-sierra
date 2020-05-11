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
Package user groups together user authentication code.
Authenticating the user, logging in, retrieving logged-in current user, etc.
*/
package user

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"github.com/swithek/sessionup"
)

const userSessionName = "userSession"

// LoginHandler handles POST /login route
type LoginHandler struct {
	userStore      Store
	sessionManager *sessionup.Manager
}

// NewLoginHandler creates a new LoginHandler
func NewLoginHandler(userStore Store, sessionManager *sessionup.Manager) *LoginHandler {
	return &LoginHandler{userStore, sessionManager}
}

// ServeHTTP handles POST /login route
func (l *LoginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	decoder := schema.NewDecoder()
	err := request.ParseForm()
	if err != nil {
		log.Println(errors.Wrap(err, "Could not parse the login form"))
		badRequest(writer, request)
		return
	}

	credentials := new(Credentials)
	decoder.ZeroEmpty(true)
	err = decoder.Decode(credentials, request.PostForm)
	if err != nil {
		log.Println(errors.Wrap(err, "Could not decode the login form into its representation"))
		http.Error(writer, "Bad Request", http.StatusBadRequest)
		return
	}

	currentUser, err := l.userStore.GetUserMatchingCredentials(credentials)
	if err != nil {
		http.Error(writer, "Forbidden", http.StatusForbidden)
		return
	}
	stringUserID := strconv.FormatUint(uint64(currentUser.ID), 10)
	err = l.sessionManager.Init(writer, request, stringUserID)
	if err != nil {
		log.Println(errors.Wrap(err, "Could not decode the user session"))
		http.Error(writer, "Internal Server Errror", http.StatusInternalServerError)
		return
	}

	http.Redirect(writer, request, "/home", http.StatusFound)
}

// badRequest replies to the request with an HTTP 400 Bad request error.
func badRequest(writer http.ResponseWriter, request *http.Request) {
	http.Error(writer, "Bad Request", http.StatusBadRequest)
}

// Credentials represents the user credentials required to log in
type Credentials struct {
	Email    string `schema:"email,required"`
	Password string `schema:"password,required"`
}

// Current represents the logged-in user
type Current struct {
	ID    uint
	Email string
}
