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
	"log"
	"net/http"
)

// ErroringHandler is a http.Handler that can return an error.
// If the error is a HTTPError, its Code and Message will be used in the response.
// If the error is not nil, it will be output as a 500 Internal Server Error.
type ErroringHandler interface {
	ServeHTTP(writer http.ResponseWriter, request *http.Request) error
}

// WrapErrors wraps ErroringHandler, catches the error returned from its ServeHTTP function
// and outputs a 500 Internal Server Error to the user with it.
func WrapErrors(next ErroringHandler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		err := next.ServeHTTP(writer, request)

		var httperr *HTTPError
		if errors.As(err, &httperr) {
			log.Println(httperr.Error())
			http.Error(writer, httperr.Message, httperr.Code)
			return
		}
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	})
}

// HTTPError is an error wrapper for errors that must be converted to HTTP errors.
// Code must be an HTTP Status code. Message is the string that will be shown to end-users.
// err is the underlying error. It is not displayed to end-users but will be logged.
type HTTPError struct {
	Code    int
	Message string
	err     error
}

// NewBadRequestError creates a new HTTPError that will be converted to a 400 Bad Request error for end-users
func NewBadRequestError(err error, message string) *HTTPError {
	return &HTTPError{http.StatusBadRequest, message, err}
}

// NewForbiddenError creates a new HTTPError that will be converted to a 403 Forbidden error
func NewForbiddenError(err error) *HTTPError {
	return &HTTPError{http.StatusForbidden, "Forbidden", err}
}

func (h *HTTPError) Unwrap() error {
	return h.err
}

func (h *HTTPError) Error() string {
	return h.Message
}
