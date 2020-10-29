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

/*
Package rest implements the REST API.
*/
package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hyzual/mike-sierra-sierra/server"
	"github.com/swithek/sessionup"
)

// Register registers a gorilla/mux Subrouter for the REST API on the given router
func Register(router *mux.Router, sessionManager *sessionup.Manager) {
	songHandler := &songHandler{}
	folderHandler := &folderHandler{}

	apiRouter := router.PathPrefix("/api/").Subrouter()
	apiRouter.Handle("/songs/{songId}", sessionManager.Auth(server.WrapErrors(songHandler)))
	apiRouter.Handle("/folders/{folderId}", sessionManager.Auth(server.WrapErrors(folderHandler)))
}

// Song represents a music file. It is distinguished by media type (audio/mp3, audio/flac, etc.)
// Most of its fields mirror tags such as ID3 tags for MP3. It is output by the REST API.
type Song struct {
	Title       string `json:"title"`       // Title is the song's title. E.g. "Know Your Enemy"
	TrackNumber uint   `json:"trackNumber"` // TrackNumber is the track number of the song. E.g. "3"
	DiskNumber  uint   `json:"diskNumber"`  // DiskNumber is the disk number of the song. E.g. "1"
	Artist      string `json:"artist"`      // Artist is the name of the main artist. E.g. "Yoko Kanno"
	Duration    uint   `json:"duration"`    // Duration is the duration of the song in seconds. E.g. "165"
	URI         string `json:"uri"`         // URI to access this song on this server. E.g. "/music/Yoko%20Kanno/1-03%20-Know%20Your%20Enemy.flac"
	Type        string `json:"type"`        // MIME type of the song E.g. "audio/flac"
}

const jsonMediaType = "application/json; charset=utf-8"

type songHandler struct {
}

func (s *songHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) error {
	response := &Song{Title: "Hello World"}
	err := json.NewEncoder(writer).Encode(response)
	if err != nil {
		return fmt.Errorf("could not encode the song %v to JSON: %w", response, err)
	}

	writer.Header().Set("Content-Type", jsonMediaType)
	writer.WriteHeader(http.StatusOK)
	return nil
}

// Folder represents a music folder. It can be any folder in the filesystem hierarchy
// such as an album, an artist folder containing many albums, a genre folder containing
// many artists, etc. It is output by the REST API.
type Folder struct {
	Name  string `json:"name"`
	Items []Song `json:"items"`
}

type folderHandler struct {
}

func (f *folderHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) error {
	songs := []Song{{Title: "Hello World"}}
	response := &Folder{Name: "Music", Items: songs}
	err := json.NewEncoder(writer).Encode(response)
	if err != nil {
		return fmt.Errorf("could not encode the folder %v to JSON: %w", response, err)
	}

	writer.Header().Set("Content-Type", jsonMediaType)
	writer.WriteHeader(http.StatusOK)
	return nil
}
