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

package user

import (
	"fmt"
	"net/http"

	"github.com/hyzual/mike-sierra-sierra/server"
	"github.com/swithek/sessionup"
)

type postSignOutHandler struct {
	sessionManager *sessionup.Manager
}

// NewSignOutPostHandler creates a new handler for POST /sign-out
func NewSignOutPostHandler(
	sessionManager *sessionup.Manager,
) http.Handler {
	return server.WrapErrors(&postSignOutHandler{sessionManager})
}

func (h *postSignOutHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) error {
	err := h.sessionManager.Revoke(request.Context(), writer)
	if err != nil {
		return fmt.Errorf("could not revoke the user session: %w", err)
	}
	http.Redirect(writer, request, "/sign-in", http.StatusFound)
	return nil
}
