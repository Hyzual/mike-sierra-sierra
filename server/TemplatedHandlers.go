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
	"fmt"
	"net/http"
)

type getLoginHandler struct {
	templateLoader TemplateLoader
	assetsResolver AssetsResolver
}

func (g *getLoginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	hashedName, err := g.assetsResolver.GetHashedName("style.css")
	if err != nil {
		http.Error(writer, fmt.Sprintf("could not resolve assets %s", err), http.StatusInternalServerError)
		return
	}
	//TODO: test those error cases
	tmpl, err := g.templateLoader.Load("login.html")
	if err != nil {
		http.Error(writer, fmt.Sprintf("could not load template %s", err), http.StatusInternalServerError)
		return
	}
	presenter := &loginPresenter{StylesheetURI: "/assets/" + hashedName}
	err = tmpl.Execute(writer, presenter)
	if err != nil {
		http.Error(writer, fmt.Sprintf("could not execute template %s", err), http.StatusInternalServerError)
		return
	}
}

type loginPresenter struct {
	StylesheetURI string
}

type homeHandler struct {
	templateLoader TemplateLoader
	assetsResolver AssetsResolver
}

func (h *homeHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	stylesheetHashedName, err := h.assetsResolver.GetHashedName("style.css")
	if err != nil {
		http.Error(writer, fmt.Sprintf("could not resolve assets %s", err), http.StatusInternalServerError)
		return
	}
	appHashedName, err := h.assetsResolver.GetHashedName("index.js")
	if err != nil {
		http.Error(writer, fmt.Sprintf("could not resolve assets %s", err), http.StatusInternalServerError)
		return
	}

	//TODO: test those error cases
	tmpl, err := h.templateLoader.Load("app.html")
	if err != nil {
		http.Error(writer, fmt.Sprintf("could not load template %s", err), http.StatusInternalServerError)
		return
	}
	presenter := &homePresenter{
		StylesheetURI: "/assets/" + stylesheetHashedName,
		AppURI:        "/assets/" + appHashedName,
	}
	err = tmpl.Execute(writer, presenter)
	if err != nil {
		http.Error(writer, fmt.Sprintf("could not execute template %s", err), http.StatusInternalServerError)
		return
	}
}

type homePresenter struct {
	StylesheetURI string // Public URI path to the stylesheet
	AppURI        string // Public URI path to the javascript app
}
