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

package user_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/hyzual/mike-sierra-sierra/server/user"
	"github.com/hyzual/mike-sierra-sierra/tests"
)

func TestLoginHandler(t *testing.T) {
	t.Run("when no request body is provided, it will return Bad Request", func(t *testing.T) {
		handler := newLoginHandlerInvalidCredentials()
		request, _ := http.NewRequest(http.MethodPost, "/login", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusBadRequest)
	})

	t.Run("when no email is provided, it will return Bad Request", func(t *testing.T) {
		handler := newLoginHandlerInvalidCredentials()
		request := newPostLoginRequest(strings.NewReader("password=welcome0"))
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusBadRequest)
	})

	t.Run("when no password is provided, it will return Bad Request", func(t *testing.T) {
		handler := newLoginHandlerInvalidCredentials()
		request := newPostLoginRequest(strings.NewReader("email=mike@example.com"))
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusBadRequest)
	})

	t.Run("when credentials don't match those from the store, it will return Forbidden", func(t *testing.T) {
		handler := newLoginHandlerInvalidCredentials()
		request := newPostLoginRequest(strings.NewReader("email=wrong.user@example.com&password=wrong_password"))
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusForbidden)
	})

	t.Run("when session cannot be opened, it will return Internal Server Error", func(t *testing.T) {
		sessionStore := &StubSessionStore{true, true}
		handler := newLoginHandlerBadSession(sessionStore)
		request := newValidPostLoginRequest()
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusInternalServerError)
	})

	t.Run("when session cannot be saved, it will return Internal Server Error", func(t *testing.T) {
		sessionStore := &StubSessionStore{false, true}
		handler := newLoginHandlerBadSession(sessionStore)
		request := newValidPostLoginRequest()
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusInternalServerError)
	})

	t.Run("when successful, POST /login will redirect to /home", func(t *testing.T) {
		handler := newValidLoginHandler()
		request := newValidPostLoginRequest()
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusFound)
		tests.AssertLocationHeaderEquals(t, response, "/home")
	})
}

func newPostLoginRequest(body io.Reader) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, "/login", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return request
}

func newValidPostLoginRequest() *http.Request {
	return newPostLoginRequest(strings.NewReader("email=mike@example.com&password=welcome0"))
}

func newLoginHandlerInvalidCredentials() *user.LoginHandler {
	dao := &StubUserDAO{false}
	sessionStore := &StubSessionStore{false, false}
	return user.NewLoginHandler(dao, sessionStore)
}

func newLoginHandlerBadSession(sessionStore sessions.Store) *user.LoginHandler {
	dao := &StubUserDAO{true}
	return user.NewLoginHandler(dao, sessionStore)
}

func newValidLoginHandler() *user.LoginHandler {
	dao := &StubUserDAO{true}
	sessionStore := &StubSessionStore{false, false}
	return user.NewLoginHandler(dao, sessionStore)
}

type StubUserDAO struct {
	isAuthenticationAccepted bool
}

func (s *StubUserDAO) GetUserMatchingCredentials(credentials *user.Credentials) (*user.Current, error) {
	if s.isAuthenticationAccepted {
		return &user.Current{ID: 1, Email: "admin@example.comn"}, nil
	}
	return nil, errors.New("Credentials do not match")
}

type StubSessionStore struct {
	shouldThrowOnGet  bool
	shouldThrowOnSave bool
}

func (s *StubSessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	if s.shouldThrowOnGet {
		return nil, errors.New("could not get session")
	}
	return sessions.NewSession(s, "testSession"), nil
}

func (s *StubSessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	return nil, errors.New("could not create new session")
}

func (s *StubSessionStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	if s.shouldThrowOnSave {
		return errors.New("could not save session")
	}
	return nil
}
