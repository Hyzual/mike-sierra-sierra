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

package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hyzual/mike-sierra-sierra/tests"
)

func TestRouter(t *testing.T) {
	router := mux.NewRouter()
	templateExecutor := newTemplateExecutorWithValidTemplate()
	assetsResolver := newValidAssetsResolver()
	userStore := newValidUserStore()
	sessionManager := tests.NewValidSessionManager(t)
	Register(router, templateExecutor, assetsResolver, userStore, sessionManager)

	t.Run("/app/suffix is handled by AppHandler", func(t *testing.T) {
		request := tests.NewAuthenticatedGetRequest(t, "/app/suffix")
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
	})
}
