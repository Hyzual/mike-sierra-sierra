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
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/schema"
	"github.com/hyzual/mike-sierra-sierra/server"
	"github.com/hyzual/mike-sierra-sierra/tests"
	"github.com/swithek/sessionup"
)

func TestGetLoginHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/login", nil)

	t.Run("when it cannot resolve assets, it will return a 500 error", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{true, ""}
		templateExecutor := newTemplateExecutorWithValidTemplate()
		handler := NewLoginGetHandler(templateExecutor, assetsResolver)

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusInternalServerError)
	})

	t.Run("when it cannot load the template, it will return a 500 error", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{false, "style.css"}
		templateExecutor := newTemplateExecutorWithInvalidTemplate()
		handler := NewLoginGetHandler(templateExecutor, assetsResolver)

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusInternalServerError)
	})

	t.Run("it will execute the template with its assets", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{false, "style.css"}
		templateExecutor := newTemplateExecutorWithValidTemplate()
		handler := NewLoginGetHandler(templateExecutor, assetsResolver)

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
	})
}

func TestPostLoginHandler(t *testing.T) {
	t.Run("when the request cannot be parsed, it will return Bad Request", func(t *testing.T) {
		handler := newLoginHandlerInvalidCredentials()
		request := httptest.NewRequest(http.MethodPost, "/login?bad-escaping-percent%", nil)
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

	t.Run("when session cannot be initialized, it will return Internal Server Error", func(t *testing.T) {
		handler := newLoginHandlerBadSession()
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
	request := httptest.NewRequest(http.MethodPost, "/login", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return request
}

func newValidPostLoginRequest() *http.Request {
	return newPostLoginRequest(strings.NewReader("email=mike@example.com&password=welcome0"))
}

func newLoginHandlerInvalidCredentials() http.Handler {
	dao := &stubDAOForLogin{false}
	sessionStore := tests.NewStubSessionStore(false)
	sessionManager := sessionup.NewManager(sessionStore)
	decoder := schema.NewDecoder()
	return NewLoginPostHandler(dao, sessionManager, decoder)
}

func newLoginHandlerBadSession() http.Handler {
	dao := &stubDAOForLogin{true}
	sessionStore := tests.NewStubSessionStore(true)
	sessionManager := sessionup.NewManager(sessionStore)
	decoder := schema.NewDecoder()
	return NewLoginPostHandler(dao, sessionManager, decoder)
}

func newValidLoginHandler() http.Handler {
	dao := &stubDAOForLogin{true}
	sessionStore := tests.NewStubSessionStore(false)
	sessionManager := sessionup.NewManager(sessionStore)
	decoder := schema.NewDecoder()
	return NewLoginPostHandler(dao, sessionManager, decoder)
}

type stubDAOForLogin struct {
	isAuthenticationAccepted bool
}

func (s *stubDAOForLogin) GetUserMatchingCredentials(_ context.Context, _ *Credentials) (*Current, error) {
	if s.isAuthenticationAccepted {
		return &Current{ID: 1, Email: "admin@example.comn"}, nil
	}
	return nil, errors.New("Credentials do not match")
}

func (s *stubDAOForLogin) SaveFirstAdministrator(_ context.Context, _ *RegistrationForm) error {
	return errors.New("This method should not have been called in tests")
}

func newTemplateExecutorWithInvalidTemplate() server.TemplateExecutor {
	return &stubTemplateExecutor{true}
}

func newTemplateExecutorWithValidTemplate() server.TemplateExecutor {
	return &stubTemplateExecutor{false}
}

type stubTemplateExecutor struct {
	shouldErrorOnLoad bool
}

func (s *stubTemplateExecutor) Load(_ io.Writer, path string, data interface{}) error {
	if s.shouldErrorOnLoad {
		return errors.New("Could not load template")
	}
	return nil
}

type stubAssetsResolver struct {
	shouldErrorOnGet bool
	filename         string
}

func (s *stubAssetsResolver) GetAssetURI(baseName string) (string, error) {
	if s.shouldErrorOnGet {
		return "", errors.New("Could not get hashed name")
	}
	return s.filename, nil
}
