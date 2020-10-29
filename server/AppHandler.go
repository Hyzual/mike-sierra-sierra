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

package server

import (
	"fmt"
	"net/http"
)

type appHandler struct {
	templateExecutor TemplateExecutor
	assetsResolver   AssetsResolver
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
	headerPresenter := &headerPresenter{"Hyzual"}
	presenter := &appPresenter{
		StylesheetURI:   styleSheetURI,
		AppURI:          scriptURI,
		HeaderPresenter: headerPresenter,
	}
	err = h.templateExecutor.Load(writer, presenter, "app.html", "sidebar.html")
	if err != nil {
		return fmt.Errorf("could not load templates %s and %s", "app.html", "sidebar.html")
	}
	return nil
}

type appPresenter struct {
	StylesheetURI   string // Public URI path to the stylesheet
	AppURI          string // Public URI path to the javascript app
	HeaderPresenter *headerPresenter
}

type headerPresenter struct {
	Username string // Username of the current logged-in user
}
