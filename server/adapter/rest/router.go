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

/*
Package rest implements the REST API.
*/
package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hyzual/mike-sierra-sierra/server/adapter/server"
	"github.com/hyzual/mike-sierra-sierra/server/domain/music"
	"github.com/swithek/sessionup"
)

// Register registers a gorilla/mux Subrouter for the REST API on the given router
func Register(router *mux.Router, sessionManager *sessionup.Manager, explorer music.MusicLibraryExplorer) {
	songHandler := &songHandler{}
	folderHandler := &folderHandler{explorer}

	apiRouter := router.PathPrefix("/api/").Subrouter()
	// All requests to the REST API must be authenticated
	apiRouter.Use(sessionManager.Auth)
	apiRouter.Handle("/songs/{songId}", server.WrapErrors(songHandler))
	apiRouter.Handle("/folders/{path:.*}", server.WrapErrors(folderHandler))
}

const jsonMediaType = "application/json; charset=utf-8"

type songHandler struct{}

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
