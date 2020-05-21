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
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hyzual/mike-sierra-sierra/tests"
)

func TestGetHome(t *testing.T) {
	request := tests.NewGetRequest(t, "/home")

	t.Run("when it cannot resolve assets, it will return an error", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{true, ""}
		templateExecutor := newTemplateExecutorWithValidTemplate()
		handler := &homeHandler{templateExecutor, assetsResolver}

		response := httptest.NewRecorder()
		err := handler.ServeHTTP(response, request)

		tests.AssertError(t, err)
	})

	t.Run("when it cannot load the template, it will return an error", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{false, "asset"}
		templateExecutor := newTemplateExecutorWithInvalidTemplate()
		handler := &homeHandler{templateExecutor, assetsResolver}

		response := httptest.NewRecorder()
		err := handler.ServeHTTP(response, request)

		tests.AssertError(t, err)
	})

	t.Run("it will execute the template with its assets", func(t *testing.T) {
		assetsResolver := &stubAssetsResolver{filename: "asset"}
		templateExecutor := newTemplateExecutorWithValidTemplate()
		handler := &homeHandler{templateExecutor, assetsResolver}

		response := httptest.NewRecorder()
		err := handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
		tests.AssertNoError(t, err)
	})
}

func newTemplateExecutorWithInvalidTemplate() TemplateExecutor {
	return &stubTemplateExecutor{true}
}

func newTemplateExecutorWithValidTemplate() TemplateExecutor {
	return &stubTemplateExecutor{false}
}

type stubTemplateExecutor struct {
	shouldErrorOnLoad bool
}

func (s *stubTemplateExecutor) Load(_ io.Writer, path string, data interface{}) error {
	if s.shouldErrorOnLoad {
		return errors.New("Could not load template")
	}
	return nil
}

type stubAssetsResolver struct {
	shouldErrorOnGet bool
	filename         string
}

func (s *stubAssetsResolver) GetAssetURI(baseName string) (string, error) {
	if s.shouldErrorOnGet {
		return "", errors.New("Could not get hashed name")
	}
	return s.filename, nil
}
