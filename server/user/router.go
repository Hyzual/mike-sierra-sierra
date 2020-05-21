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
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/hyzual/mike-sierra-sierra/server"
	"github.com/swithek/sessionup"
)

// Register registers the routes on the given gorilla/mux router.
func Register(
	router *mux.Router,
	templateExecutor server.TemplateExecutor,
	assetsResolver server.AssetsResolver,
	userStore Store,
	sessionManager *sessionup.Manager,
	decoder *schema.Decoder,
) {
	getLoginHandler := NewLoginGetHandler(templateExecutor, assetsResolver)
	postLoginHandler := NewLoginPostHandler(userStore, sessionManager, decoder)
	getFirstTimeRegistrationHandler := NewFirstTimeRegistrationGetHandler(templateExecutor, assetsResolver)
	postFirstTimeRegistrationHandler := NewFirstTimeRegistrationPostHandler(userStore, decoder)

	router.Handle("/first-time-registration", getFirstTimeRegistrationHandler).Methods(http.MethodGet)
	router.Handle("/first-time-registration", postFirstTimeRegistrationHandler).Methods(http.MethodPost)
	router.Handle("/login", getLoginHandler).Methods(http.MethodGet)
	router.Handle("/login", postLoginHandler).Methods(http.MethodPost)
}
