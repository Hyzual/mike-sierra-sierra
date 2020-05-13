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

type getLoginHandler struct {
	templateLoader TemplateLoader
	assetsResolver AssetsResolver
}

func (g *getLoginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) error {
	hashedName, err := g.assetsResolver.GetHashedName("style.css")
	if err != nil {
		return errors.Wrapf(err, "could not resolve assets %s", "style.css")
	}
	tmpl, err := g.templateLoader.Load("login.html")
	if err != nil {
		return errors.Wrapf(err, "could not load template %s", "login.html")
	}
	presenter := &loginPresenter{StylesheetURI: "/assets/" + hashedName}
	return tmpl.Execute(writer, presenter)
}

type loginPresenter struct {
	StylesheetURI string
}

type homeHandler struct {
	templateLoader TemplateLoader
	assetsResolver AssetsResolver
}

func (h *homeHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) error {
	stylesheetHashedName, err := h.assetsResolver.GetHashedName("style.css")
	if err != nil {
		return errors.Wrapf(err, "could not resolve asset %s", "style.css")
	}
	appHashedName, err := h.assetsResolver.GetHashedName("index.js")
	if err != nil {
		return errors.Wrapf(err, "could not resolve asset %s", "index.js")
	}
	tmpl, err := h.templateLoader.Load("app.html")
	if err != nil {
		return errors.Wrapf(err, "could not load template %s", "app.html")
	}
	headerPresenter := &headerPresenter{"Hyzual"}
	presenter := &homePresenter{
		StylesheetURI:   "/assets/" + stylesheetHashedName,
		AppURI:          "/assets/" + appHashedName,
		HeaderPresenter: headerPresenter,
	}
	return tmpl.Execute(writer, presenter)
}

type homePresenter struct {
	StylesheetURI   string // Public URI path to the stylesheet
	AppURI          string // Public URI path to the javascript app
	HeaderPresenter *headerPresenter
}

type headerPresenter struct {
	Username string // Username of the current logged-in user
}
