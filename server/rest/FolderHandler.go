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
