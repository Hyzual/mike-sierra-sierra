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
	"fmt"
	"net/http"
)

// TopLevelFolders represents the list of top-level music folders. It is a struct wrapping a slice
// instead of a slice to avoid producing a JSON array susceptible to JSON hijacking.
type TopLevelFolders struct {
	Folders []SubFolder `json:"folders"`
}

// SubFolder represents a music folder that is a child of a Folder. We do not expose its items
// yet, another query is needed to expose them.
type SubFolder struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type topFolderHandler struct{}

func (h *topFolderHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) error {
	folders := []SubFolder{
		{ID: 1, Name: "Ghost in the Shell - Stand Alone Complex OST 3"},
		{ID: 2, Name: "Call To Power 2"},
		{ID: 3, Name: "Civilization: Call To Power"},
		{ID: 4, Name: "Medieval II Total War"},
		{ID: 5, Name: "Age of Empires Definitive Edition (Original Soundtrack)"},
		{ID: 6, Name: "Stellaris Digital Soundtrack"},
		{ID: 7, Name: "WarCraft III: Reign of Chaos [Ripped]"},
		{ID: 8, Name: "Zeus: Master of Olympus"},
		{ID: 9, Name: "Il était une fois... l'Homme"},
		{ID: 10, Name: "Ghost in the Shell - Stand Alone Complex : Solid State Society"},
		{ID: 11, Name: "Final Fantasy X OST"},
		{ID: 12, Name: "Video Games Music"},
		{ID: 13, Name: "Trance"},
		{ID: 14, Name: "Starcraft OST"},
	}
	response := &TopLevelFolders{folders}
	err := json.NewEncoder(writer).Encode(response)
	if err != nil {
		return fmt.Errorf("could not encode the top-level folders to JSON: %w", err)
	}

	writer.Header().Set("Content-Type", jsonMediaType)
	writer.WriteHeader(http.StatusOK)
	return nil
}

// Folder represents a music folder. It can be any folder in the filesystem hierarchy
// such as an album, an artist folder containing many albums, a genre folder containing
// many artists, etc. It is output by the REST API.
type Folder struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Items []Song `json:"items"`
}

type folderHandler struct{}

func (f *folderHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) error {
	songs := []Song{{Title: "Hello World"}}
	response := &Folder{ID: 1, Name: "Music", Items: songs}
	err := json.NewEncoder(writer).Encode(response)
	if err != nil {
		return fmt.Errorf("could not encode the folder %v to JSON: %w", response, err)
	}

	writer.Header().Set("Content-Type", jsonMediaType)
	writer.WriteHeader(http.StatusOK)
	return nil
}
