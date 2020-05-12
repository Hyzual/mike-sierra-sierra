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
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hyzual/mike-sierra-sierra/tests"
)

func TestGetLogin(t *testing.T) {
	request := tests.NewGetRequest(t, "/login")

	t.Run("when it cannot resolve assets, it will return a 500 error", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{true, ""}
		templateLoader := newTemplateLoaderWithValidTemplate()
		musicServer := newMusicServerWithDeps(assetsResolver, templateLoader)

		response := httptest.NewRecorder()
		musicServer.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusInternalServerError)
	})

	t.Run("when it cannot load the template, it will return a 500 error", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{false, "style.css"}
		templateLoader := newTemplateLoaderWithInvalidTemplate()
		musicServer := newMusicServerWithDeps(assetsResolver, templateLoader)

		response := httptest.NewRecorder()
		musicServer.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusInternalServerError)
	})

	t.Run("when it cannot execute the template, it will return a 500 error", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{false, "style.css"}
		templateLoader := newTemplateLoaderWithTemplateExecError()
		musicServer := newMusicServerWithDeps(assetsResolver, templateLoader)

		response := httptest.NewRecorder()
		musicServer.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusInternalServerError)
	})

	t.Run("it will execute the template with its assets", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{false, "style.css"}
		templateLoader := newTemplateLoaderWithValidTemplate()
		musicServer := newMusicServerWithDeps(assetsResolver, templateLoader)

		response := httptest.NewRecorder()
		musicServer.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
	})
}

func TestGetHome(t *testing.T) {
	request := tests.NewGetRequest(t, "/home")

	t.Run("when it cannot resolve assets, it will return an error", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{true, ""}
		templateLoader := newTemplateLoaderWithValidTemplate()
		handler := &homeHandler{templateLoader, assetsResolver}

		response := httptest.NewRecorder()

		err := handler.ServeHTTP(response, request)

		tests.AssertError(t, err)
	})

	t.Run("when it cannot load the template, it will return an error", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{false, "asset"}
		templateLoader := newTemplateLoaderWithInvalidTemplate()
		handler := &homeHandler{templateLoader, assetsResolver}

		response := httptest.NewRecorder()
		err := handler.ServeHTTP(response, request)

		tests.AssertError(t, err)
	})

	t.Run("when it cannot execute the template, it will return an error", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{false, "asset"}
		templateLoader := newTemplateLoaderWithTemplateExecError()
		handler := &homeHandler{templateLoader, assetsResolver}

		response := httptest.NewRecorder()
		err := handler.ServeHTTP(response, request)

		tests.AssertError(t, err)
	})

	t.Run("it will execute the template with its assets", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{filename: "asset"}
		templateLoader := newTemplateLoaderWithValidTemplate()
		handler := &homeHandler{templateLoader, assetsResolver}

		response := httptest.NewRecorder()

		err := handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
		tests.AssertNoError(t, err)
	})
}

func newMusicServerWithDeps(assetsResolver AssetsResolver, templateLoader TemplateLoader) *MusicServer {
	sessionManager := newSessionManager()
	return New(
		sessionManager,
		nil,
		templateLoader,
		nil,
		assetsResolver,
		nil,
	)
}

func newTemplateLoaderWithInvalidTemplate() TemplateLoader {
	return &stubTemplateLoader{true, nil}
}

func newTemplateLoaderWithTemplateExecError() TemplateLoader {
	stubTemplate, _ := template.New("stub").Parse(`<html>{{template "Unknown"}}</html>`)
	return &stubTemplateLoader{false, stubTemplate}
}

func newTemplateLoaderWithValidTemplate() TemplateLoader {
	stubTemplate, err := template.New("stub").Parse(`<html></html>`)
	fmt.Fprintln(os.Stdout, err)
	return &stubTemplateLoader{false, stubTemplate}
}

type stubTemplateLoader struct {
	shouldErrorOnLoad bool
	tmpl              *template.Template
}

func (s *stubTemplateLoader) Load(path string) (*template.Template, error) {
	if s.shouldErrorOnLoad {
		return nil, errors.New("Could not load template")
	}
	return s.tmpl, nil
}

type stubAssetsResolver struct {
	shouldErrorOnGet bool
	filename         string
}

func (s *stubAssetsResolver) GetHashedName(baseName string) (string, error) {
	if s.shouldErrorOnGet {
		return "", errors.New("Could not get hashed name")
	}
	return s.filename, nil
}
