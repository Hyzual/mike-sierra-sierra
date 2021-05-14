/*
 *   Copyright (C) 2020-2021  Joris MASSON
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

package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hyzual/mike-sierra-sierra/tests"
)

func TestRouter(t *testing.T) {
	router := mux.NewRouter()
	sessionManager := tests.NewValidSessionManager(t)
	assetsLoader := &stubPathJoiner{filename: ""}
	musicLoader := &stubPathJoiner{filename: ""}
	Register(router, sessionManager, assetsLoader, musicLoader)

	t.Run("/unknown returns 404", func(t *testing.T) {
		request := tests.NewGetRequest(t, "/unknown")
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusNotFound)
	})

	t.Run("/ redirects to /app", func(t *testing.T) {
		request := tests.NewGetRequest(t, "/")
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusFound)
		tests.AssertLocationHeaderEquals(t, response, "/app")
	})
}

func TestGetAssets(t *testing.T) {
	tempFile, removeTempFile := createTempFile(t)
	defer removeTempFile()
	assetsLoader := &stubPathJoiner{tempFile.Name()}
	handler := &assetsHandler{assetsLoader}

	t.Run("returns OK for a path leading to a file", func(t *testing.T) {
		request := tests.NewGetRequest(t, "/assets/style.css")
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
	})

	t.Run("returns NotFound for a path leading to /assets directory", func(t *testing.T) {
		request := tests.NewGetRequest(t, "/assets/")
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusNotFound)
	})
}

func TestGetMusic(t *testing.T) {
	tempFile, removeTempFile := createTempFile(t)
	defer removeTempFile()
	musicLoader := &stubPathJoiner{tempFile.Name()}
	handler := &musicHandler{musicLoader}

	t.Run("returns OK for a path leading to a music file", func(t *testing.T) {
		request := tests.NewGetRequest(t, "/music/album/amazing-song.mp3")
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
	})
}

func TestUnauthorized(t *testing.T) {
	handler := HandleUnauthorized(errors.New("Error"))
	request := tests.NewGetRequest(t, "/app")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	tests.AssertStatusEquals(t, response.Code, http.StatusFound)
	tests.AssertLocationHeaderEquals(t, response, "/sign-in")
}

type stubPathJoiner struct {
	filename string
}

func (s *stubPathJoiner) Join(relativePath string) string {
	return s.filename
}

func createTempFile(t *testing.T) (*os.File, func()) {
	t.Helper()

	tempFile, err := os.CreateTemp("", "db")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	removeFile := func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}

	return tempFile, removeFile
}
