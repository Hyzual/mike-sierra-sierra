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

	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
)

const userSessionName = "userSession"

// LoginHandler handles POST /login route
type LoginHandler struct {
	userStore    Store
	sessionStore sessions.Store
}

// NewLoginHandler create a new LoginHandler
func NewLoginHandler(userStore Store, sessionStore sessions.Store) *LoginHandler {
	return &LoginHandler{userStore, sessionStore}
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
	session, err := l.sessionStore.Get(request, userSessionName)
	if err != nil {
		log.Println(errors.Wrap(err, "Could not decode the user session"))
		http.Error(writer, "Internal Server Errror", http.StatusInternalServerError)
		return
	}
	session.Values["userID"] = currentUser.ID
	session.Options = &sessions.Options{
		Path:     "/",
		Secure:   true,
		MaxAge:   86400 * 7, // session should last one week
		SameSite: http.SameSiteStrictMode,
	}
	err = session.Save(request, writer)
	if err != nil {
		log.Println(errors.Wrap(err, "Could not save the user session"))
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
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
