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

/*
 Package app groups together app routes. Displaying the frontend app
 to the user with its assets.
*/
package app

import (
	"github.com/gorilla/mux"
	"github.com/hyzual/mike-sierra-sierra/server/adapter"
	"github.com/hyzual/mike-sierra-sierra/server/adapter/server"
	"github.com/hyzual/mike-sierra-sierra/server/adapter/server/user"
	"github.com/swithek/sessionup"
)

// Register registers the routes on the given gorilla/mux router.
func Register(
	router *mux.Router,
	templateExecutor adapter.TemplateExecutor,
	assetsResolver adapter.AssetsResolver,
	userStore user.Store,
	sessionManager *sessionup.Manager,
) {
	appHandler := sessionManager.Auth(
		server.WrapErrors(&appHandler{templateExecutor, assetsResolver, userStore}),
	)
	router.PathPrefix("/app").Handler(appHandler)
}
