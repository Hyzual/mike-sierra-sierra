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

package tests

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/swithek/sessionup"
)

// NewGetRequest simplifies creating a new GET request
func NewGetRequest(t *testing.T, url string) *http.Request {
	t.Helper()
	return httptest.NewRequest(http.MethodGet, url, nil)
}

// NewAuthenticatedGetRequest creates a new GET request to url with a cookie named "id"
func NewAuthenticatedGetRequest(t *testing.T, url string) *http.Request {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, url, nil)
	request.AddCookie(&http.Cookie{Name: "id"})
	return request
}

// AssertStatusEquals verifies that the request status code matches expectation
func AssertStatusEquals(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

// AssertLocationHeaderEquals verifies that the response's Location header matches expectation
func AssertLocationHeaderEquals(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	got := response.Header().Get("Location")
	if got != want {
		t.Errorf("Location Header did not match expected route, got %s, want %s", got, want)
	}
}

// AssertContentTypeHeaderEquals verifies that the response's Content-Type header matches expectation
func AssertContentTypeHeaderEquals(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	got := response.Header().Get("Content-Type")
	if got != want {
		t.Errorf("Content-Type Header did not match expected type, got %s, want %s", got, want)
	}
}

// AssertNoError verifies that the error is nil
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("did not expect an error, got one %v", err)
	}
}

// AssertError verifies that the error is not nil
func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected an error but did not get one")
	}
}

// NewValidSessionManager creates a new valid session manager with cookie named "id"
func NewValidSessionManager(t *testing.T) *sessionup.Manager {
	t.Helper()
	store := NewValidSessionStore(t)
	return sessionup.NewManager(store, sessionup.CookieName("id"))
}

// stubSessionStore mocks a sessionup Store
type stubSessionStore struct {
	shouldThrowOnCreate bool
	shouldThrowOnDelete bool
}

// NewValidSessionStore creates a new session store that never throws
func NewValidSessionStore(t *testing.T) sessionup.Store {
	t.Helper()
	return &stubSessionStore{false, false}
}

// NewStoreWithErrorOnCreate creates a new session store that throws on Create()
func NewStoreWithErrorOnCreate(t *testing.T) sessionup.Store {
	t.Helper()
	return &stubSessionStore{true, false}
}

// NewSessionStoreWithErrorOnDelete creates a new session store that throws on Delete()
func NewSessionStoreWithErrorOnDelete(t *testing.T) sessionup.Store {
	t.Helper()
	return &stubSessionStore{false, true}
}

// Create mocks sessionup Store's method.
func (s *stubSessionStore) Create(ctx context.Context, session sessionup.Session) error {
	if s.shouldThrowOnCreate {
		return errors.New("Could not create session")
	}
	return nil
}

// FetchByID mocks sessionup Store's method
func (s *stubSessionStore) FetchByID(ctx context.Context, id string) (sessionup.Session, bool, error) {
	return sessionup.Session{ID: id}, true, nil
}

// FetchByUserKey mocks sessionup Store's method
func (s *stubSessionStore) FetchByUserKey(ctx context.Context, key string) ([]sessionup.Session, error) {
	return nil, errors.New("This method is not supposed to be call in the tests")
}

// DeleteByID mocks sessionup Store's method
func (s *stubSessionStore) DeleteByID(ctx context.Context, id string) error {
	if s.shouldThrowOnDelete {
		return errors.New("Could not delete the session")
	}
	return nil
}

// DeleteByUserKey mocks sessionup Store's method
func (s *stubSessionStore) DeleteByUserKey(ctx context.Context, key string, expID ...string) error {
	return errors.New("This method is not supposed to be call in the tests")
}
