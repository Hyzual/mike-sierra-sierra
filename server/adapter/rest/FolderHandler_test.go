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
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hyzual/mike-sierra-sierra/server/domain/music"
	"github.com/hyzual/mike-sierra-sierra/tests"
)

func TestGetFolder(t *testing.T) {
	response := httptest.NewRecorder()

	t.Run(`given no path, it will return the JSON representation of the contents of the root music library folder`, func(t *testing.T) {
		request := tests.NewGetRequest(t, "/api/folders/")
		explorer := newValidLibraryExplorer(t)
		handler := &folderHandler{explorer}

		err := handler.ServeHTTP(response, request)
		tests.AssertNoError(t, err)

		var got Folder
		err = json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("Unable to parse response from server %q into Folder, %v", response.Body, err)
		}
		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
		tests.AssertContentTypeHeaderEquals(t, response, jsonMediaType)
	})

	t.Run(`given a path, it will return the JSON representation of the folder's contents at that path`, func(t *testing.T) {
		request := newGetRequestWithPathVar(t, "Sub Folder")
		explorer := newValidLibraryExplorer(t)
		handler := &folderHandler{explorer}

		err := handler.ServeHTTP(response, request)
		tests.AssertNoError(t, err)

		var got Folder
		err = json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("Unable to parse response from server %q into Folder, %v", response.Body, err)
		}
		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
		tests.AssertContentTypeHeaderEquals(t, response, jsonMediaType)
	})

	t.Run(`when folders or songs are nil slices, it will return an empty JSON array instead of null`, func(t *testing.T) {
		request := newGetRequestWithPathVar(t, "Sub Folder")
		explorer := newLibraryExplorerReturnsEmpty(t)
		handler := &folderHandler{explorer}

		err := handler.ServeHTTP(response, request)
		tests.AssertNoError(t, err)

		var got Folder
		err = json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("Unable to parse response from server %q into Folder, %v", response.Body, err)
		}
		if got.Folders == nil {
			t.Errorf("Folders should be empty array, not nil")
		}
		if got.Songs == nil {
			t.Errorf("Songs should be empty array, not nil")
		}
	})

	t.Run("it will return an error if the given folder cannot be read", func(t *testing.T) {
		request := newGetRequestWithPathVar(t, "unknown")
		explorer := newLibraryExplorerWithError(t)
		handler := &folderHandler{explorer}

		err := handler.ServeHTTP(response, request)
		tests.AssertError(t, err)
	})
}

func newGetRequestWithPathVar(t *testing.T, pathName string) *http.Request {
	request := tests.NewGetRequest(t, "/api/folders/"+url.PathEscape(pathName))
	vars := make(map[string]string)
	vars["path"] = pathName
	mux.SetURLVars(request, vars)
	return request
}

func newValidLibraryExplorer(t *testing.T) music.MusicLibraryExplorer {
	t.Helper()
	return &libraryExplorerStub{true, false}
}

func newLibraryExplorerReturnsEmpty(t *testing.T) music.MusicLibraryExplorer {
	t.Helper()
	return &libraryExplorerStub{true, true}
}

func newLibraryExplorerWithError(t *testing.T) music.MusicLibraryExplorer {
	t.Helper()
	return &libraryExplorerStub{false, false}
}

type libraryExplorerStub struct {
	isValid    bool
	returnsNil bool
}

func (e *libraryExplorerStub) ListContents(pathName string) ([]music.SubFolder, []music.Song, error) {
	if !e.isValid {
		return nil, nil, errors.New("This error should be expected in tests")
	}
	if e.returnsNil {
		folders := []music.SubFolder{}
		songs := []music.Song{}
		return folders, songs, nil
	}
	folders := []music.SubFolder{
		{Name: "satisfied", Path: "Sub Folder/satisfied"},
		{Name: "indicate", Path: "Sub Folder/indicate"},
	}
	songs := []music.Song{
		{Title: "Medicine Worry", Path: "Sub Folder/Medicine Worry.mp3"},
		{Title: "He Wall", Path: "Sub Folder/He Wall.ogg"},
	}
	return folders, songs, nil
}
