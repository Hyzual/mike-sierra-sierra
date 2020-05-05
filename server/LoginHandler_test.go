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

package server_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hyzual/mike-sierra-sierra/server"
)

func TestGetLogin(t *testing.T) {
	handler := buildLoginHandler(nil)
	request := newGetRequest(t, "/login")
	response := httptest.NewRecorder()

	handler.GetHandler(response, request)

	assertStatusEquals(t, response.Code, http.StatusOK)
}

func TestPostLogin(t *testing.T) {
	userDAO := StubUserDAO{false}
	handler := buildLoginHandler(&userDAO)

	t.Run("when no request body is provided, it will return Bad Request", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/login", nil)
		response := httptest.NewRecorder()

		handler.PostHandler(response, request)

		assertStatusEquals(t, response.Code, http.StatusBadRequest)
	})

	t.Run("when no email is provided, it will return Bad Request", func(t *testing.T) {
		request := newPostLoginRequest(t, strings.NewReader("password=welcome0"))
		response := httptest.NewRecorder()

		handler.PostHandler(response, request)

		assertStatusEquals(t, response.Code, http.StatusBadRequest)
	})

	t.Run("when no password is provided, it will return Bad Request", func(t *testing.T) {
		request := newPostLoginRequest(t, strings.NewReader("email=mike@example.com"))
		response := httptest.NewRecorder()

		handler.PostHandler(response, request)

		assertStatusEquals(t, response.Code, http.StatusBadRequest)
	})

	t.Run("when credentials don't match those from the store, it will return Forbidden", func(t *testing.T) {
		request := newPostLoginRequest(t, strings.NewReader("email=wrong.user@example.com&password=wrong_password"))
		response := httptest.NewRecorder()

		handler.PostHandler(response, request)

		assertStatusEquals(t, response.Code, http.StatusForbidden)
	})

	t.Run("when successful, POST /login will redirect to /home", func(t *testing.T) {
		userDAO = StubUserDAO{true}
		handler = buildLoginHandler(&userDAO)
		request := newPostLoginRequest(t, strings.NewReader("email=mike@example.com&password=welcome0"))
		response := httptest.NewRecorder()

		handler.PostHandler(response, request)

		assertStatusEquals(t, response.Code, http.StatusFound)
		assertLocationHeaderEquals(t, response, "/home")
	})
}

func newPostLoginRequest(t *testing.T, body io.Reader) *http.Request {
	t.Helper()
	request, _ := http.NewRequest(http.MethodPost, "/login", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return request
}

func buildLoginHandler(userStore server.UserStore) *server.LoginHandler {
	pathJoiner := newTestPathJoiner()
	return server.NewLoginHandler(pathJoiner, userStore)
}

type StubUserDAO struct {
	isAuthenticationAccepted bool
}

func (s *StubUserDAO) VerifyCredentialsMatch(server.LoginFormRepresentation) error {
	if s.isAuthenticationAccepted {
		return nil
	}
	return errors.New("Credentials do not match")
}
