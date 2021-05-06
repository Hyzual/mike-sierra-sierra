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

package server

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hyzual/mike-sierra-sierra/tests"
)

func TestUnkwnownRoute(t *testing.T) {
	musicServer := newMusicServer()
	t.Run("/unknown returns 404", func(t *testing.T) {
		request := tests.NewGetRequest(t, "/unknown")
		response := httptest.NewRecorder()
		musicServer.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusNotFound)
	})
}

func TestGetRoot(t *testing.T) {
	musicServer := newMusicServer()

	t.Run("/ redirects to /app", func(t *testing.T) {
		request := tests.NewGetRequest(t, "/")
		response := httptest.NewRecorder()
		musicServer.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusFound)
		tests.AssertLocationHeaderEquals(t, response, "/app")
	})
}

func TestGetAssets(t *testing.T) {
	tempFile, removeTempFile := createTempFile(t)
	defer removeTempFile()
	musicServer := newMusicServerWithAsset(tempFile.Name())

	t.Run("returns OK for a path leading to a file", func(t *testing.T) {
		request := tests.NewGetRequest(t, "/assets/style.css")
		response := httptest.NewRecorder()

		musicServer.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
	})

	t.Run("returns NotFound for a path leading to /assets directory", func(t *testing.T) {
		request := tests.NewGetRequest(t, "/assets/")
		response := httptest.NewRecorder()

		musicServer.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusNotFound)
	})
}

func TestGetMusic(t *testing.T) {
	tempFile, removeTempFile := createTempFile(t)
	defer removeTempFile()
	musicServer := newMusicServerWithMusic(tempFile.Name())

	t.Run("returns OK for a path leading to a music file", func(t *testing.T) {
		request := tests.NewGetRequest(t, "/music/album/amazing-song.mp3")
		response := httptest.NewRecorder()

		musicServer.ServeHTTP(response, request)

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

func newMusicServer() *MusicServer {
	router := mux.NewRouter()
	assetsLoader := &stubPathJoiner{filename: ""}
	return New(router, assetsLoader, nil)
}

func newMusicServerWithAsset(filename string) *MusicServer {
	router := mux.NewRouter()
	assetsLoader := &stubPathJoiner{filename}
	return New(router, assetsLoader, nil)
}

func newMusicServerWithMusic(filename string) *MusicServer {
	router := mux.NewRouter()
	musicLoader := &stubPathJoiner{filename}
	return New(router, nil, musicLoader)
}

type stubPathJoiner struct {
	filename string
}

func (s *stubPathJoiner) Join(relativePath string) string {
	return s.filename
}

func createTempFile(t *testing.T) (*os.File, func()) {
	t.Helper()

	tempFile, err := ioutil.TempFile("", "db")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	removeFile := func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}

	return tempFile, removeFile
}
