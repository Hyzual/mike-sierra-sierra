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
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hyzual/mike-sierra-sierra/tests"
	"github.com/swithek/sessionup"
)

func TestPostSignOutHandler(t *testing.T) {
	t.Run("when the session manager cannot revoke the user session, it will return Internal Server Error", func(t *testing.T) {
		handler := newSignOutHandlerBadSession()
		request := newPostSignOutRequest()
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusInternalServerError)
	})

	t.Run("when succesful, POST /sign-out will redirect to /sign-in", func(t *testing.T) {
		handler := newValidSignOutHandler()
		request := newPostSignOutRequest()
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusFound)
		tests.AssertLocationHeaderEquals(t, response, "/sign-in")
	})
}

func newSignOutHandlerBadSession() http.Handler {
	sessionStore := tests.NewStubSessionStore(false, true)
	sessionManager := sessionup.NewManager(sessionStore)
	return NewSignOutPostHandler(sessionManager)
}

func newValidSignOutHandler() http.Handler {
	sessionStore := tests.NewStubSessionStore(false, false)
	sessionManager := sessionup.NewManager(sessionStore)
	return NewSignOutPostHandler(sessionManager)
}

func newPostSignOutRequest() *http.Request {
	request := httptest.NewRequest(http.MethodPost, "/sign-out", strings.NewReader(""))
	session := sessionup.Session{
		ID: "id",
	}
	return request.WithContext(sessionup.NewContext(request.Context(), session))
}
