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

package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hyzual/mike-sierra-sierra/tests"
)

func TestWrapErrors(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/endpoint", nil)

	t.Run("when the wrapped handler returns nil, it does nothing", func(t *testing.T) {
		handler := WrapErrors(&stubErroringHandler{nil})
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
	})

	t.Run(`when the wrapped handler returns an HTTP Error,
		it will respond with the error's message and code`, func(t *testing.T) {
		err := NewBadRequestError(errors.New("Bad format"), "Bad format received in /endpoint")
		handler := WrapErrors(&stubErroringHandler{err})
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusBadRequest)
	})

	t.Run(`when the wrapped handler returns a non-HTTP error,
		it will respond with a 500 Internal Server Error`, func(t *testing.T) {
		err := errors.New("Could not connect to database")
		handler := WrapErrors(&stubErroringHandler{err})
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusInternalServerError)
	})
}

type stubErroringHandler struct {
	err error
}

func (s *stubErroringHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) error {
	return s.err
}
