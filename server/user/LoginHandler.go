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

package user

import (
	"net/http"
	"strconv"

	"github.com/gorilla/schema"
	"github.com/hyzual/mike-sierra-sierra/server"
	"github.com/pkg/errors"
	"github.com/swithek/sessionup"
)

const userSessionName = "userSession"

// NewLoginGetHandler creates a new handler for GET /login
func NewLoginGetHandler(
	te server.TemplateExecutor,
	ar server.AssetsResolver,
) http.Handler {
	return server.WrapErrors(
		&getLoginHandler{te, ar},
	)
}

type getLoginHandler struct {
	templateExecutor server.TemplateExecutor
	assetsResolver   server.AssetsResolver
}

func (h *getLoginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) error {
	hashedName, err := h.assetsResolver.GetAssetURI("style.css")
	if err != nil {
		return errors.Wrapf(err, "could not resolve assets %s", "style.css")
	}
	presenter := &loginPresenter{StylesheetURI: hashedName}
	err = h.templateExecutor.Load(writer, "login.html", presenter)
	if err != nil {
		return errors.Wrapf(err, "could not load template %s", "login.html")
	}
	return nil
}

type loginPresenter struct {
	StylesheetURI string
}

type postLoginHandler struct {
	userStore      Store
	sessionManager *sessionup.Manager
	decoder        *schema.Decoder
}

// NewLoginPostHandler creates a new handler for POST /login
func NewLoginPostHandler(
	userStore Store,
	sessionManager *sessionup.Manager,
	decoder *schema.Decoder,
) http.Handler {
	return server.WrapErrors(
		&postLoginHandler{userStore, sessionManager, decoder},
	)
}

func (h *postLoginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) error {
	err := request.ParseForm()
	if err != nil {
		return server.NewBadRequestError(err, "Could not parse the login form")
	}

	credentials := new(Credentials)
	err = h.decoder.Decode(credentials, request.PostForm)
	if err != nil {
		return server.NewBadRequestError(err, "Could not decode the login form into its representation")
	}

	currentUser, err := h.userStore.GetUserMatchingCredentials(request.Context(), credentials)
	if err != nil {
		return server.NewForbiddenError(errors.New("Invalid credentials"))
	}
	stringUserID := strconv.FormatUint(uint64(currentUser.ID), 10)
	err = h.sessionManager.Init(writer, request, stringUserID)
	if err != nil {
		return errors.Wrap(err, "Could not decode the user session")
	}

	http.Redirect(writer, request, "/home", http.StatusFound)
	return nil
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
