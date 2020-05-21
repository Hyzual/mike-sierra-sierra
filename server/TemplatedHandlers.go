/*
 *   Copyright (c) 2020 Joris MASSON

 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.

 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.

 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package server

import (
	"net/http"

	"github.com/pkg/errors"
)

type homeHandler struct {
	templateExecutor TemplateExecutor
	assetsResolver   AssetsResolver
}

func (h *homeHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) error {
	styleSheetURI, err := h.assetsResolver.GetAssetURI("style.css")
	if err != nil {
		return errors.Wrapf(err, "could not resolve asset %s", "style.css")
	}
	scriptURI, err := h.assetsResolver.GetAssetURI("index.js")
	if err != nil {
		return errors.Wrapf(err, "could not resolve asset %s", "index.js")
	}
	headerPresenter := &headerPresenter{"Hyzual"}
	presenter := &homePresenter{
		StylesheetURI:   styleSheetURI,
		AppURI:          scriptURI,
		HeaderPresenter: headerPresenter,
	}
	err = h.templateExecutor.Load(writer, "app.html", presenter)
	if err != nil {
		return errors.Wrapf(err, "could not load template %s", "app.html")
	}
	return nil
}

type homePresenter struct {
	StylesheetURI   string // Public URI path to the stylesheet
	AppURI          string // Public URI path to the javascript app
	HeaderPresenter *headerPresenter
}

type headerPresenter struct {
	Username string // Username of the current logged-in user
}
