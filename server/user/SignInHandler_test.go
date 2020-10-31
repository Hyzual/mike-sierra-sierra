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

func TestGetSignInHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/sign-in", nil)

	t.Run("when it cannot resolve assets, it will return a 500 error", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{true, ""}
		templateExecutor := newTemplateExecutorWithValidTemplate()
		handler := NewSignInGetHandler(templateExecutor, assetsResolver)

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusInternalServerError)
	})

	t.Run("when it cannot load the template, it will return a 500 error", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{false, "style.css"}
		templateExecutor := newTemplateExecutorWithInvalidTemplate()
		handler := NewSignInGetHandler(templateExecutor, assetsResolver)

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusInternalServerError)
	})

	t.Run("it will execute the template with its assets", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{false, "style.css"}
		templateExecutor := newTemplateExecutorWithValidTemplate()
		handler := NewSignInGetHandler(templateExecutor, assetsResolver)

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
	})
}

func TestPostSigninHandler(t *testing.T) {
	t.Run("when the request cannot be parsed, it will return Bad Request", func(t *testing.T) {
		handler := newSignInHandlerInvalidCredentials()
		request := httptest.NewRequest(http.MethodPost, "/sign-in?bad-escaping-percent%", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusBadRequest)
	})

	t.Run("when no email is provided, it will return Bad Request", func(t *testing.T) {
		handler := newSignInHandlerInvalidCredentials()
		request := newPostSigninRequest(strings.NewReader("password=welcome0"))
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusBadRequest)
	})

	t.Run("when no password is provided, it will return Bad Request", func(t *testing.T) {
		handler := newSignInHandlerInvalidCredentials()
		request := newPostSigninRequest(strings.NewReader("email=mike@example.com"))
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusBadRequest)
	})

	t.Run("when the email address doesn't match any from the store, it will return Forbidden", func(t *testing.T) {
		handler := newSignInHandlerInvalidCredentials()
		request := newPostSigninRequest(strings.NewReader("email=wrong.user@example.com&password=wrong_password"))
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusForbidden)
	})

	t.Run("when the password does not match the password Hash from the store, it will return Forbidden", func(t *testing.T) {
		handler := newSignInHandlerBadSession()
		request := newPostSigninRequest(strings.NewReader("email=mike@example.com&password=wrong_password"))
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusForbidden)
	})

	t.Run("when session cannot be initialized, it will return Internal Server Error", func(t *testing.T) {
		handler := newSignInHandlerBadSession()
		request := newValidPostSigninRequest()
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusInternalServerError)
	})

	t.Run("when successful, POST /sign-in will redirect to /app", func(t *testing.T) {
		handler := newValidSignInHandler()
		request := newValidPostSigninRequest()
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusFound)
		tests.AssertLocationHeaderEquals(t, response, "/app")
	})
}

func newPostSigninRequest(body io.Reader) *http.Request {
	request := httptest.NewRequest(http.MethodPost, "/sign-in", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return request
}

func newValidPostSigninRequest() *http.Request {
	return newPostSigninRequest(strings.NewReader("email=mike@example.com&password=welcome0"))
}

func newSignInHandlerInvalidCredentials() http.Handler {
	dao := &stubDAOForSignIn{false}
	sessionStore := tests.NewStubSessionStore(false, false)
	sessionManager := sessionup.NewManager(sessionStore)
	decoder := schema.NewDecoder()
	return NewSignInPostHandler(dao, sessionManager, decoder)
}

func newSignInHandlerBadSession() http.Handler {
	dao := &stubDAOForSignIn{true}
	sessionStore := tests.NewStubSessionStore(true, false)
	sessionManager := sessionup.NewManager(sessionStore)
	decoder := schema.NewDecoder()
	return NewSignInPostHandler(dao, sessionManager, decoder)
}

func newValidSignInHandler() http.Handler {
	dao := &stubDAOForSignIn{true}
	sessionStore := tests.NewStubSessionStore(false, false)
	sessionManager := sessionup.NewManager(sessionStore)
	decoder := schema.NewDecoder()
	return NewSignInPostHandler(dao, sessionManager, decoder)
}

type stubDAOForSignIn struct {
	isAuthenticationAccepted bool
}

// Corresponds to "welcome0" password with work = 4
var testPasswordHash = []byte{36, 50, 97, 36, 48, 52, 36, 101, 74, 54, 110, 79, 74, 118, 115, 100, 68, 117, 86, 50, 76, 116, 65, 69, 55, 76, 76, 109, 46, 80, 85, 78, 98, 89, 120, 122, 72, 104, 117, 99, 50, 112, 72, 116, 100, 55, 122, 114, 76, 73, 106, 117, 46, 119, 100, 50, 87, 52, 118, 109}

func (s *stubDAOForSignIn) GetUserMatchingEmail(_ context.Context, _ string) (*PossibleMatch, error) {
	if s.isAuthenticationAccepted {
		return &PossibleMatch{ID: 1, Email: "admin@example.com", PasswordHash: testPasswordHash}, nil
	}
	return nil, errors.New("Credentials do not match")
}

func (s *stubDAOForSignIn) SaveFirstAdministrator(_ context.Context, _ *Registration) error {
	return errors.New("This method should not have been called in tests")
}

func (s *stubDAOForSignIn) GetUserMatchingSession(_ context.Context) (*Current, error) {
	return nil, errors.New("This method should not have been called in tests")
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

func (s *stubTemplateExecutor) Load(_ io.Writer, data interface{}, templatePaths ...string) error {
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
