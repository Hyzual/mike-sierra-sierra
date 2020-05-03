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

package server_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/hyzual/mike-sierra-sierra/server"
)

func TestGetRoot(t *testing.T) {
	musicServer := buildMusicServerWithCorrectRootDir()

	t.Run("/ redirects to /home", func(t *testing.T) {
		request := newGetRequest(t, "/")
		response := httptest.NewRecorder()
		musicServer.ServeHTTP(response, request)

		assertStatusEquals(t, response.Code, http.StatusFound)
		assertLocationHeaderEquals(t, response, "/home")
	})

	t.Run("/unknown route will return NotFound", func(t *testing.T) {
		request := newGetRequest(t, "/unknown")
		response := httptest.NewRecorder()
		musicServer.ServeHTTP(response, request)

		assertStatusEquals(t, response.Code, http.StatusNotFound)
	})
}

func TestGetHome(t *testing.T) {
	musicServer := buildMusicServerWithCorrectRootDir()
	request := newGetRequest(t, "/home")
	response := httptest.NewRecorder()

	musicServer.ServeHTTP(response, request)

	assertStatusEquals(t, response.Code, http.StatusOK)
}

func TestGetAssets(t *testing.T) {
	musicServer := buildMusicServerWithCorrectRootDir()

	t.Run("returns OK for a path leading to a file", func(t *testing.T) {
		request := newGetRequest(t, "/assets/style.css")
		response := httptest.NewRecorder()

		musicServer.ServeHTTP(response, request)

		assertStatusEquals(t, response.Code, http.StatusOK)
	})

	t.Run("forbids listing directories", func(t *testing.T) {
		request := newGetRequest(t, "/assets/")
		response := httptest.NewRecorder()

		musicServer.ServeHTTP(response, request)

		assertStatusEquals(t, response.Code, http.StatusForbidden)
	})
}

func TestGetLogin(t *testing.T) {
	musicServer := buildMusicServerWithCorrectRootDir()
	request, _ := http.NewRequest(http.MethodGet, "/login", nil)
	response := httptest.NewRecorder()

	musicServer.ServeHTTP(response, request)

	assertStatusEquals(t, response.Code, http.StatusOK)
}

func buildMusicServerWithCorrectRootDir() *server.MusicServer {
	wd, _ := os.Getwd() // will be <projectRoot>/server
	rootDir := filepath.Join(wd, "..")
	return server.New(rootDir)
}

func newGetRequest(t *testing.T, url string) *http.Request {
	t.Helper()
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return req
}

func assertStatusEquals(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertLocationHeaderEquals(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	got := response.Header().Get("Location")
	if got != want {
		t.Errorf("Location Header did not match expected route, got %s, want %s", got, want)
	}
}
