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
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hyzual/mike-sierra-sierra/server"
	"github.com/hyzual/mike-sierra-sierra/server/user"
	"github.com/hyzual/mike-sierra-sierra/tests"
)

func TestGetApp(t *testing.T) {
	request := tests.NewGetRequest(t, "/app")

	t.Run("when it cannot resolve assets, it will return an error", func(t *testing.T) {
		assetsResolver := newInvalidAssetsResolver()
		templateExecutor := newTemplateExecutorWithValidTemplate()
		userStore := newValidUserStore()
		handler := &appHandler{templateExecutor, assetsResolver, userStore}

		response := httptest.NewRecorder()
		err := handler.ServeHTTP(response, request)

		tests.AssertError(t, err)
	})

	t.Run("when it cannot load the template, it will return an error", func(t *testing.T) {
		assetsResolver := newValidAssetsResolver()
		templateExecutor := newTemplateExecutorWithInvalidTemplate()
		userStore := newValidUserStore()
		handler := &appHandler{templateExecutor, assetsResolver, userStore}

		response := httptest.NewRecorder()
		err := handler.ServeHTTP(response, request)

		tests.AssertError(t, err)
	})

	t.Run("when it cannot retrieve the current user, it will return an error", func(t *testing.T) {
		assetsResolver := newValidAssetsResolver()
		templateExecutor := newTemplateExecutorWithValidTemplate()
		userStore := newInvalidUserStore()
		handler := &appHandler{templateExecutor, assetsResolver, userStore}

		response := httptest.NewRecorder()
		err := handler.ServeHTTP(response, request)

		tests.AssertError(t, err)
	})

	t.Run("it will execute the template with its assets", func(t *testing.T) {
		assetsResolver := newValidAssetsResolver()
		templateExecutor := newTemplateExecutorWithValidTemplate()
		userStore := newValidUserStore()
		handler := &appHandler{templateExecutor, assetsResolver, userStore}

		response := httptest.NewRecorder()
		err := handler.ServeHTTP(response, request)

		tests.AssertStatusEquals(t, response.Code, http.StatusOK)
		tests.AssertNoError(t, err)
	})
}

func newTemplateExecutorWithInvalidTemplate() server.TemplateExecutor {
	return &stubTemplateExecutor{true}
}

func newTemplateExecutorWithValidTemplate() server.TemplateExecutor {
	return &stubTemplateExecutor{false}
}

type stubTemplateExecutor struct {
	shouldErrorOnLoad bool
}

func (s *stubTemplateExecutor) Load(_ io.Writer, data interface{}, templatePaths ...string) error {
	if s.shouldErrorOnLoad {
		return errors.New("Could not load template")
	}
	return nil
}

func newValidAssetsResolver() server.AssetsResolver {
	return &stubAssetsResolver{filename: "asset"}
}

func newInvalidAssetsResolver() server.AssetsResolver {
	return &stubAssetsResolver{true, ""}
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

func newValidUserStore() user.Store {
	currentUser := &user.CurrentUser{ID: 27, Email: "testuser@example.com", Username: "Test User"}
	return &stubDAOForApp{false, currentUser}
}

func newInvalidUserStore() user.Store {
	return &stubDAOForApp{true, nil}
}

type stubDAOForApp struct {
	shouldErrorOnGet bool
	currentUser      *user.CurrentUser
}

func (s *stubDAOForApp) GetUserMatchingSession(_ context.Context) (*user.CurrentUser, error) {
	if s.shouldErrorOnGet {
		return nil, errors.New("Could not get current user")
	}
	return s.currentUser, nil
}

func (s *stubDAOForApp) GetUserMatchingEmail(_ context.Context, _ string) (*user.PossibleMatch, error) {
	return nil, errors.New("This method should not have been called in tests")
}

func (s *stubDAOForApp) SaveFirstAdministrator(_ context.Context, _ *user.Registration) error {
	return errors.New("This method should not have been called in tests")
}

func TestGenerateGravatarHash(t *testing.T) {
	t.Run("it generates a md5 hash of the trimmed and lowercased email from the current user", func(t *testing.T) {
		currentUser := &user.CurrentUser{Email: "Valid.Email@example.com"}
		expected := "cef5ba9259f7619f438306c020cda589"

		actual := generateGravatarHash(currentUser)
		if actual != expected {
			t.Errorf("expected hash %s to be %s", actual, expected)
		}
	})

	t.Run("it returns zeroes when there is no email", func(t *testing.T) {
		currentUser := &user.CurrentUser{}
		expected := "00000000000000000000000000000000"
		actual := generateGravatarHash(currentUser)
		if actual != expected {
			t.Errorf("expected hash %s to be %s", actual, expected)
		}
	})
}
