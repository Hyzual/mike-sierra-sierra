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

package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hyzual/mike-sierra-sierra/tests"
)

func TestGetSong(t *testing.T) {
	handler := &songHandler{}
	t.Run("given a song ID, it will return JSON containing Hello World", func(t *testing.T) {
		request := tests.NewGetRequest(t, "/api/songs/15")
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		var got Song
		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("could not decode the response body from server %q into a Song, '%v'", response.Body, err)
		}

		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
		tests.AssertContentTypeHeaderEquals(t, response, jsonMediaType)
	})
}

func TestGetFolder(t *testing.T) {
	handler := &folderHandler{}
	t.Run(`when folder id 0 is given,
		it will return the representation of the top-level music folder`, func(t *testing.T) {
		request := tests.NewGetRequest(t, "/api/folders/0")
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		var got Folder
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Folder, '%v'", response.Body, err)
		}
		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
		tests.AssertContentTypeHeaderEquals(t, response, jsonMediaType)
	})
}
