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
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return req
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

// StubSessionStore mocks a sessionup Store
type StubSessionStore struct {
	shouldThrowOnCreate bool
}

// NewStubSessionStore creates a new Stub store
func NewStubSessionStore(shouldThrowOnCreate bool) *StubSessionStore {
	return &StubSessionStore{shouldThrowOnCreate}
}

// Create mocks sessionup Store's method. It throws an error when StubSessionStore.shouldThrowOnCreate is true
func (s *StubSessionStore) Create(ctx context.Context, session sessionup.Session) error {
	if s.shouldThrowOnCreate {
		return errors.New("Could not create session")
	}
	return nil
}

// FetchByID mocks sessionup Store's method
func (s *StubSessionStore) FetchByID(ctx context.Context, id string) (sessionup.Session, bool, error) {
	return sessionup.Session{}, false, errors.New("This method is not supposed to be call in the tests")
}

// FetchByUserKey mocks sessionup Store's method
func (s *StubSessionStore) FetchByUserKey(ctx context.Context, key string) ([]sessionup.Session, error) {
	return nil, errors.New("This method is not supposed to be call in the tests")
}

// DeleteByID mocks sessionup Store's method
func (s *StubSessionStore) DeleteByID(ctx context.Context, id string) error {
	return errors.New("This method is not supposed to be call in the tests")
}

// DeleteByUserKey mocks sessionup Store's method
func (s *StubSessionStore) DeleteByUserKey(ctx context.Context, key string, expID ...string) error {
	return errors.New("This method is not supposed to be call in the tests")
}
