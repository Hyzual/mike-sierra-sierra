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
	"github.com/hyzual/mike-sierra-sierra/tests"
)

func TestGetFirstTimeRegistrationHandler(t *testing.T) {
	request := tests.NewGetRequest(t, "/first-time-registration")

	t.Run("when it cannot resolve assets, it will return a 500 error", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{true, ""}
		templateExecutor := newTemplateExecutorWithValidTemplate()
		handler := NewFirstTimeRegistrationGetHandler(templateExecutor, assetsResolver)

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusInternalServerError)
	})

	t.Run("when it cannot load the template, it will return a 500 error", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{false, "style.css"}
		templateExecutor := newTemplateExecutorWithInvalidTemplate()
		handler := NewFirstTimeRegistrationGetHandler(templateExecutor, assetsResolver)

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusInternalServerError)
	})

	t.Run("it will execute the template with its assets", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{false, "style.css"}
		templateExecutor := newTemplateExecutorWithValidTemplate()
		handler := NewFirstTimeRegistrationGetHandler(templateExecutor, assetsResolver)

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
	})
}

func TestPostFirstTimeRegistrationHandler(t *testing.T) {
	t.Run("when the request cannot be parsed, it will return Bad Request", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/first-time-registration?bad-escaping-percent%", nil)
		response := httptest.NewRecorder()
		handler := newFirstTimeRegistrationHandler()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusBadRequest)
	})

	t.Run("when no email is provided, it will return Bad Request", func(t *testing.T) {
		request := newPostFirstRegistrationRequest(strings.NewReader("password=welcome0"))
		response := httptest.NewRecorder()
		handler := newFirstTimeRegistrationHandler()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusBadRequest)
	})

	t.Run("when no password is provided it will return Bad Request", func(t *testing.T) {
		request := newPostFirstRegistrationRequest(strings.NewReader("email=mike@example.com"))
		response := httptest.NewRecorder()
		handler := newFirstTimeRegistrationHandler()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusBadRequest)
	})

	t.Run("when no username is provided, it will return Bad Request", func(t *testing.T) {
		request := newPostFirstRegistrationRequest(strings.NewReader("email=mike@example.com&password=welcome0"))
		response := httptest.NewRecorder()
		handler := newFirstTimeRegistrationHandler()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusBadRequest)
	})

	t.Run(`when the password is longer than 64 characters, it will return Bad Request`, func(t *testing.T) {
		request := newPostFirstRegistrationRequest(
			strings.NewReader("email=mike@example.com&username=mike&password=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		)
		response := httptest.NewRecorder()
		handler := newFirstTimeRegistrationHandler()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusBadRequest)
	})

	t.Run("when it can't register the admin in the database, it will return Internal Server Error", func(t *testing.T) {
		request := newValidPostFirstRegistrationRequest()
		response := httptest.NewRecorder()
		handler := newFirstTimeRegistrationHandlerWithDBError()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusInternalServerError)
	})

	t.Run("when successful, POST /first-time-registration will redirect to /sign-in", func(t *testing.T) {
		request := newValidPostFirstRegistrationRequest()
		response := httptest.NewRecorder()
		handler := newFirstTimeRegistrationHandler()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusFound)
		tests.AssertLocationHeaderEquals(t, response, "/sign-in")
	})
}

func newFirstTimeRegistrationHandler() http.Handler {
	dao := &stubDAOForRegistration{false}
	decoder := schema.NewDecoder()
	return NewFirstTimeRegistrationPostHandler(dao, decoder)
}

func newFirstTimeRegistrationHandlerWithDBError() http.Handler {
	dao := &stubDAOForRegistration{true}
	decoder := schema.NewDecoder()
	return NewFirstTimeRegistrationPostHandler(dao, decoder)
}

func newPostFirstRegistrationRequest(body io.Reader) *http.Request {
	request := httptest.NewRequest(http.MethodPost, "/first-time-registration", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return request
}

func newValidPostFirstRegistrationRequest() *http.Request {
	return newPostFirstRegistrationRequest(strings.NewReader("email=mike@example.com&password=welcome0&username=Mike"))
}

type stubDAOForRegistration struct {
	hasError bool
}

func (s *stubDAOForRegistration) SaveFirstAdministrator(_ context.Context, registration *Registration) error {
	if s.hasError {
		return errors.New("Could not register first administrator")
	}
	return nil
}

func (s *stubDAOForRegistration) GetUserMatchingEmail(_ context.Context, _ string) (*PossibleMatch, error) {
	return nil, errors.New("This method should not have been called in tests")
}

func (s *stubDAOForRegistration) GetUserMatchingSession(_ context.Context) (*CurrentUser, error) {
	return nil, errors.New("This method should not have been called in tests")
}
