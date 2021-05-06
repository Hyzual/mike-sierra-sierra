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

package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hyzual/mike-sierra-sierra/tests"
)

func TestRouter(t *testing.T) {
	router := mux.NewRouter()
	sessionManager := tests.NewValidSessionManager(t)
	explorer := newValidLibraryExplorer(t)
	Register(router, sessionManager, explorer)

	t.Run("/api/folders/path is handled by FolderHandler", func(t *testing.T) {
		request := tests.NewAuthenticatedGetRequest(t, "/api/folders/path")
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
	})

	t.Run("/api/songs/1 is handled by SongHandler", func(t *testing.T) {
		request := tests.NewAuthenticatedGetRequest(t, "/api/songs/1")
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
	})
}

func TestGetSong(t *testing.T) {
	handler := &songHandler{}
	t.Run("given a song ID, it will return JSON containing Hello World", func(t *testing.T) {
		request := tests.NewGetRequest(t, "/api/songs/15")
		response := httptest.NewRecorder()

		err := handler.ServeHTTP(response, request)
		tests.AssertNoError(t, err)

		var got Song
		err = json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("could not decode the response body from server %q into a Song, '%v'", response.Body, err)
		}

		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
		tests.AssertContentTypeHeaderEquals(t, response, jsonMediaType)
	})
}
