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

package server_test

import (
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/swithek/sessionup"

	"github.com/hyzual/mike-sierra-sierra/server"
	"github.com/hyzual/mike-sierra-sierra/tests"
)

func TestGetRoot(t *testing.T) {
	musicServer := newMusicServer()

	t.Run("/ redirects to /home", func(t *testing.T) {
		request := tests.NewGetRequest(t, "/")
		response := httptest.NewRecorder()
		musicServer.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusFound)
		tests.AssertLocationHeaderEquals(t, response, "/home")
	})

	t.Run("/unknown route will return NotFound", func(t *testing.T) {
		request := tests.NewGetRequest(t, "/unknown")
		response := httptest.NewRecorder()
		musicServer.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusNotFound)
	})
}

func TestGetHome(t *testing.T) {
	templateLoader := &StubTemplateLoader{}
	handler := server.NewHomeHandler(templateLoader)

	request := tests.NewGetRequest(t, "/home")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	tests.AssertStatusEquals(t, response.Code, http.StatusOK)
}

func TestGetLogin(t *testing.T) {
	musicServer := newMusicServer()
	request := tests.NewGetRequest(t, "/login")
	response := httptest.NewRecorder()

	musicServer.ServeHTTP(response, request)

	tests.AssertStatusEquals(t, response.Code, http.StatusOK)
}

func TestGetAssets(t *testing.T) {
	tempFile, removeTempFile := createTempFile(t)
	defer removeTempFile()
	musicServer := newMusicServerWithAsset(tempFile.Name())

	t.Run("returns OK for a path leading to a file", func(t *testing.T) {
		request := tests.NewGetRequest(t, "/assets/style.css")
		response := httptest.NewRecorder()

		musicServer.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
	})

	t.Run("returns NotFound for a path leading to /assets directory", func(t *testing.T) {
		request := tests.NewGetRequest(t, "/assets/")
		response := httptest.NewRecorder()

		musicServer.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusNotFound)
	})
}

func TestGetMusic(t *testing.T) {
	tempFile, removeTempFile := createTempFile(t)
	defer removeTempFile()
	musicServer := newMusicServerWithMusic(tempFile.Name())

	t.Run("returns OK for a path leading to a music file", func(t *testing.T) {
		request := tests.NewGetRequest(t, "/music/album/amazing-song.mp3")
		response := httptest.NewRecorder()

		musicServer.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
	})
}

func TestUnauthorized(t *testing.T) {
	handler := server.HandleUnauthorized(errors.New("Error"))
	request := tests.NewGetRequest(t, "/home")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	tests.AssertStatusEquals(t, response.Code, http.StatusFound)
	tests.AssertLocationHeaderEquals(t, response, "/login")
}

func newSessionManager() *sessionup.Manager {
	sessionStore := tests.NewStubSessionStore(false)
	return sessionup.NewManager(
		sessionStore,
		sessionup.CookieName("id"),
		sessionup.Reject(server.HandleUnauthorized),
	)
}

func newMusicServer() *server.MusicServer {
	sessionManager := newSessionManager()
	assetsIncluder := &StubAssetsIncluder{filename: ""}
	templateLoader := &StubTemplateLoader{}
	return server.New(sessionManager, assetsIncluder, templateLoader, nil, nil)
}

func newMusicServerWithAsset(filename string) *server.MusicServer {
	sessionManager := newSessionManager()
	assetsIncluder := &StubAssetsIncluder{filename}
	templateLoader := &StubTemplateLoader{}
	return server.New(sessionManager, assetsIncluder, templateLoader, nil, nil)
}

func newMusicServerWithMusic(filename string) *server.MusicServer {
	sessionManager := newSessionManager()
	musicLoader := &StubMusicLoader{filename}
	return server.New(sessionManager, nil, nil, musicLoader, nil)
}

type StubAssetsIncluder struct {
	filename string
}

func (s *StubAssetsIncluder) Join(relativePath string) string {
	return s.filename
}

type StubTemplateLoader struct{}

func (s *StubTemplateLoader) Load(path string) (*template.Template, error) {
	stubTemplate, _ := template.New("stub").Parse("<html></html>")
	return stubTemplate, nil
}

type StubMusicLoader struct {
	filename string
}

func (s *StubMusicLoader) Join(relativePath string) string {
	return s.filename
}

func createTempFile(t *testing.T) (*os.File, func()) {
	t.Helper()

	tempFile, err := ioutil.TempFile("", "db")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	removeFile := func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}

	return tempFile, removeFile
}
