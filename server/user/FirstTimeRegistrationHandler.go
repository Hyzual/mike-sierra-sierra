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

	"github.com/gorilla/schema"
	"github.com/hyzual/mike-sierra-sierra/server"
	"github.com/pkg/errors"
)

// NewFirstTimeRegistrationGetHandler creates a new handler for GET /first-time-registration
func NewFirstTimeRegistrationGetHandler(
	te server.TemplateExecutor,
	ar server.AssetsResolver,
) http.Handler {
	return server.WrapErrors(
		&getFirstTimeRegistrationHandler{te, ar},
	)
}

type getFirstTimeRegistrationHandler struct {
	templateExecutor server.TemplateExecutor
	assetsResolver   server.AssetsResolver
}

func (h *getFirstTimeRegistrationHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) error {
	styleSheetURI, err := h.assetsResolver.GetAssetURI("style.css")
	if err != nil {
		return errors.Wrapf(err, "could not resolve asset %s", "style.css")
	}
	presenter := &firstTimeRegistrationPresenter{StylesheetURI: styleSheetURI}
	err = h.templateExecutor.Load(writer, "first-time-registration.html", presenter)
	if err != nil {
		return errors.Wrapf(err, "could not load template %s", "first-time-registration.html")
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
	err = h.userStore.SaveFirstAdministrator(request.Context(), form)
	if err != nil {
		return errors.Wrap(err, "Error while saving the first administrator account")
	}

	http.Redirect(writer, request, "/login", http.StatusFound)
	return nil
}

// RegistrationForm represents the registration information provided by users
// to create their account
type RegistrationForm struct {
	Credentials
	Username string `schema:"username,required"`
}
