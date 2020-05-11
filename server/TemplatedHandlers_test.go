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
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hyzual/mike-sierra-sierra/tests"
)

func TestGetLogin(t *testing.T) {
	musicServer := newMusicServerWithManifest("style.css")
	request := tests.NewGetRequest(t, "/login")
	response := httptest.NewRecorder()

	musicServer.ServeHTTP(response, request)

	tests.AssertStatusEquals(t, response.Code, http.StatusOK)
}

func TestGetHome(t *testing.T) {
	templateLoader := &stubTemplateLoader{}
	assetsResolver := &stubAssetsResolver{filename: "style.css"}
	handler := &homeHandler{templateLoader, assetsResolver}

	request := tests.NewGetRequest(t, "/home")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	tests.AssertStatusEquals(t, response.Code, http.StatusOK)
}

func newMusicServerWithManifest(hashedName string) *MusicServer {
	sessionManager := newSessionManager()
	templateLoader := &stubTemplateLoader{}
	assetsResolver := &stubAssetsResolver{hashedName}
	return New(sessionManager, nil, templateLoader, nil, assetsResolver, nil)
}

type stubTemplateLoader struct{}

func (s *stubTemplateLoader) Load(path string) (*template.Template, error) {
	stubTemplate, _ := template.New("stub").Parse("<html></html>")
	return stubTemplate, nil
}

type stubAssetsResolver struct {
	filename string
}

func (s *stubAssetsResolver) GetHashedName(baseName string) (string, error) {
	return s.filename, nil
}
