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
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hyzual/mike-sierra-sierra/server/domain/music"
)

// FolderContents represents the contents of a folder. It can contain SubFolders and Songs
type FolderContents struct {
	Folders []SubFolder `json:"folders"`
	Songs   []Song      `json:"songs"`
}

// SubFolder represents a music folder that is a child of a Folder. We do not expose its items
// yet, another HTTP request is needed to expose them.
type SubFolder struct {
	Name string `json:"name"` // Basename of the folder
	Path string `json:"path"` // Path from the root music folder. For example "Yoko%20Kanno"
}

func fromSubFolder(source music.SubFolder) SubFolder {
	return SubFolder{
		Name: source.Name,
		Path: source.Path,
	}
}

// Song represents a music file. It is distinguished by media type (audio/mp3, audio/flac, etc.)
// Most of its fields mirror tags such as ID3 tags for MP3. It is output by the REST API.
type Song struct {
	Title string `json:"title"` // Title is the song's title. E.g. "Know Your Enemy"
	// TrackNumber uint   `json:"trackNumber"` // TrackNumber is the track number of the song. E.g. "3"
	// DiskNumber  uint   `json:"diskNumber"`  // DiskNumber is the disk number of the song. E.g. "1"
	// Artist      string `json:"artist"`      // Artist is the name of the main artist. E.g. "Yoko Kanno"
	// Duration    uint   `json:"duration"`    // Duration is the duration of the song in seconds. E.g. "165"
	URI string `json:"uri"` // URI to access this song on this server. E.g. "/music/Yoko%20Kanno/1-03%20-Know%20Your%20Enemy.flac"
	// Type        string `json:"type"`        // MIME type of the song E.g. "audio/flac"
}

func fromSong(source music.Song) Song {
	return Song{
		Title: source.Title,
		URI:   source.URI,
	}
}

func mapIntoRepresentations(contentFolders []music.SubFolder, contentSongs []music.Song) FolderContents {
	folders := make([]SubFolder, 0) // Init slice at zero, otherwise nil slice results in "null" JSON instead of []
	songs := make([]Song, 0)
	for _, folder := range contentFolders {
		folders = append(folders, fromSubFolder(folder))
	}
	for _, song := range contentSongs {
		songs = append(songs, fromSong(song))
	}
	return FolderContents{folders, songs}
}

// Folder represents a music folder. It can be any folder in the filesystem hierarchy
// such as an album, an artist folder containing many albums, a genre folder containing
// many artists, etc. It is output by the REST API.
type Folder struct {
	Folders []SubFolder `json:"folders"`
	Songs   []Song      `json:"songs"`
}

type folderHandler struct {
	explorer music.MusicLibraryExplorer
}

func (h *folderHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) error {
	vars := mux.Vars(request)
	folderPath := vars["path"]
	if folderPath == "" {
		folderPath = "." // List the contents of the root music library folder
	}
	folders, songs, err := h.explorer.ListContents(folderPath)
	if err != nil {
		return fmt.Errorf("error while retrieving the contents of the folder at path %v: %w", folderPath, err)
	}
	contents := mapIntoRepresentations(folders, songs)
	response := Folder{Folders: contents.Folders, Songs: contents.Songs}
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		return fmt.Errorf("could not encode the folder %v to JSON: %w", response, err)
	}

	writer.Header().Set("Content-Type", jsonMediaType)
	writer.WriteHeader(http.StatusOK)
	return nil
}
