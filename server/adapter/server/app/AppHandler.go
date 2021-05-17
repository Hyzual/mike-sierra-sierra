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

package app

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"crypto/md5" //nolint gosec //md5 is required for Gravatar and is not used for sensitive crypto here

	"github.com/hyzual/mike-sierra-sierra/server/adapter"
	"github.com/hyzual/mike-sierra-sierra/server/adapter/server/user"
)

type appHandler struct {
	templateExecutor adapter.TemplateExecutor
	assetsResolver   adapter.AssetsResolver
	userStore        user.Store
}

func (h *appHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) error {
	styleSheetURI, err := h.assetsResolver.GetAssetURI("style.css")
	if err != nil {
		return fmt.Errorf("could not resolve asset %s: %w", "style.css", err)
	}
	scriptURI, err := h.assetsResolver.GetAssetURI("index.js")
	if err != nil {
		return fmt.Errorf("could not resolve asset %s: %w", "index.js", err)
	}
	currentUser, err := h.userStore.GetUserMatchingSession(request.Context())
	if err != nil {
		return fmt.Errorf("could not retrieve the current user: %w", err)
	}
	headerPresenter := &headerPresenter{currentUser.Username, generateGravatarHash(currentUser)}
	presenter := &appPresenter{
		StylesheetURI:   styleSheetURI,
		AppURI:          scriptURI,
		HeaderPresenter: headerPresenter,
	}
	err = h.templateExecutor.Load(writer, presenter, "app.html")
	if err != nil {
		return fmt.Errorf("could not load template %s", "app.html")
	}
	return nil
}

// See https://gravatar.com/site/implement/hash/
func generateGravatarHash(currentUser *user.Current) string {
	hasher := md5.New() //nolint gosec //md5 is required for Gravatar and is not used for sensitive crypto here
	email := strings.ToLower(strings.TrimSpace(currentUser.Email))
	if email == "" {
		return "00000000000000000000000000000000"
	}
	_, err := hasher.Write([]byte(email))
	if err != nil {
		return "00000000000000000000000000000000"
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

type appPresenter struct {
	StylesheetURI   string // Public URI path to the stylesheet
	AppURI          string // Public URI path to the javascript app
	HeaderPresenter *headerPresenter
}

type headerPresenter struct {
	Username     string // Username of the current logged-in user
	GravatarHash string // Gravatar hash of the current logged-in user email.
}
