/*
 *   Copyright (C) 2020  Joris MASSON
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

package user

import (
	"net/http"
	"strconv"

	"github.com/gorilla/schema"
	"github.com/hyzual/mike-sierra-sierra/server"
	"github.com/pkg/errors"
	"github.com/swithek/sessionup"
	"golang.org/x/crypto/bcrypt"
)

const userSessionName = "userSession"

// NewSignInGetHandler creates a new handler for GET /sign-in
func NewSignInGetHandler(
	te server.TemplateExecutor,
	ar server.AssetsResolver,
) http.Handler {
	return server.WrapErrors(
		&getSignInHandler{te, ar},
	)
}

type getSignInHandler struct {
	templateExecutor server.TemplateExecutor
	assetsResolver   server.AssetsResolver
}

func (h *getSignInHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) error {
	hashedName, err := h.assetsResolver.GetAssetURI("style.css")
	if err != nil {
		return errors.Wrapf(err, "could not resolve assets %s", "style.css")
	}
	presenter := &signInPresenter{StylesheetURI: hashedName}
	err = h.templateExecutor.Load(writer, "sign-in.html", presenter)
	if err != nil {
		return errors.Wrapf(err, "could not load template %s", "sign-in.html")
	}
	return nil
}

type signInPresenter struct {
	StylesheetURI string
}

type postSignInHandler struct {
	userStore      Store
	sessionManager *sessionup.Manager
	decoder        *schema.Decoder
}

// NewSignInPostHandler creates a new handler for POST /sign-in
func NewSignInPostHandler(
	userStore Store,
	sessionManager *sessionup.Manager,
	decoder *schema.Decoder,
) http.Handler {
	return server.WrapErrors(
		&postSignInHandler{userStore, sessionManager, decoder},
	)
}

func (h *postSignInHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) error {
	err := request.ParseForm()
	if err != nil {
		return server.NewBadRequestError(err, "Could not parse the sign-in form")
	}

	signInForm := new(SignInForm)
	err = h.decoder.Decode(signInForm, request.PostForm)
	if err != nil {
		return server.NewBadRequestError(err, "Could not decode the sign-in form into its representation")
	}

	possibleUser, err := h.userStore.GetUserMatchingEmail(request.Context(), signInForm.Email)
	if err != nil {
		return server.NewForbiddenError(errors.New("Invalid credentials"))
	}
	err = bcrypt.CompareHashAndPassword(possibleUser.PasswordHash, []byte(signInForm.Password))
	if err != nil {
		return server.NewForbiddenError(errors.New("Invalid credentials"))
	}

	stringUserID := strconv.FormatUint(uint64(possibleUser.ID), 10)
	err = h.sessionManager.Init(writer, request, stringUserID)
	if err != nil {
		return errors.Wrap(err, "Could not decode the user session")
	}

	http.Redirect(writer, request, "/home", http.StatusFound)
	return nil
}

// SignInForm represents the credentials provided by users to sign in
type SignInForm struct {
	Email    string `schema:"email,required"`
	Password string `schema:"password,required"`
}