/*
 *   Copyright (C) 2020-2021  Joris MASSON
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
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/hyzual/mike-sierra-sierra/server/adapter"
	"github.com/hyzual/mike-sierra-sierra/server/adapter/server"
	"golang.org/x/crypto/bcrypt"
)

// Passwords are limited to 64 characters because bcrypt is limited to 72 characters
// but we don't want to reveal we're using bcrypt.
// See https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html#maximum-password-lengths
const (
	maximumPasswordLength = 64
	bcryptWork            = 12
)

// NewFirstTimeRegistrationGetHandler creates a new handler for GET /first-time-registration
func NewFirstTimeRegistrationGetHandler(
	te adapter.TemplateExecutor,
	ar adapter.AssetsResolver,
) http.Handler {
	return server.WrapErrors(
		&getFirstTimeRegistrationHandler{te, ar},
	)
}

type getFirstTimeRegistrationHandler struct {
	templateExecutor adapter.TemplateExecutor
	assetsResolver   adapter.AssetsResolver
}

func (h *getFirstTimeRegistrationHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) error {
	styleSheetURI, err := h.assetsResolver.GetAssetURI("style.css")
	if err != nil {
		return fmt.Errorf("could not resolve asset %s: %w", "style.css", err)
	}
	presenter := &firstTimeRegistrationPresenter{StylesheetURI: styleSheetURI}
	err = h.templateExecutor.Load(writer, presenter, "first-time-registration.html")
	if err != nil {
		return fmt.Errorf("could not load template %s: %w", "first-time-registration.html", err)
	}
	return nil
}

type firstTimeRegistrationPresenter struct {
	StylesheetURI string // Public URI path to the stylesheet
}

// NewFirstTimeRegistrationPostHandler creates a new handler for POST /first-time-registration
func NewFirstTimeRegistrationPostHandler(
	us Store,
	de *schema.Decoder,
) http.Handler {
	return server.WrapErrors(
		&postFirstTimeRegistrationHandler{us, de},
	)
}

// postFirstTimeRegistrationHandler handles POST /first-time-registration route
type postFirstTimeRegistrationHandler struct {
	userStore Store
	decoder   *schema.Decoder
}

func (h *postFirstTimeRegistrationHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) error {
	err := request.ParseForm()
	if err != nil {
		return server.NewBadRequestError(err, "Could not parse the first-time registration form")
	}
	form := new(RegistrationForm)
	err = h.decoder.Decode(form, request.PostForm)
	if err != nil {
		return server.NewBadRequestError(err, "Could not decode the first-time registration form into its representation")
	}
	if len([]rune(form.Password)) > maximumPasswordLength {
		return server.NewBadRequestError(err, "Password cannot be longer than 64 characters")
	}

	passwordhash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcryptWork)
	if err != nil {
		return fmt.Errorf("error while hashing the password: %w", err)
	}

	registration := &Registration{
		Email:        form.Email,
		PasswordHash: passwordhash,
		Username:     form.Username,
	}

	err = h.userStore.SaveFirstAdministrator(request.Context(), registration)
	if err != nil {
		return fmt.Errorf("error while saving the first administrator account: %w", err)
	}

	http.Redirect(writer, request, "/sign-in", http.StatusFound)
	return nil
}

// RegistrationForm represents the registration information provided by users
// to create their account
type RegistrationForm struct {
	Email    string `schema:"email,required"`
	Password string `schema:"password,required"`
	Username string `schema:"username,required"`
}
