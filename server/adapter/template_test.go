/*
 *   Copyright (C) 2020-2021  Joris MASSON
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

package adapter_test

import (
	"os"
	"strings"
	"testing"

	"github.com/hyzual/mike-sierra-sierra/server/adapter"
	"github.com/hyzual/mike-sierra-sierra/tests"
)

func TestTemplateExecutor(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not get the current working directory, '%v'", err)
	}
	loader := adapter.NewTemplateExecutor(basePath)

	t.Run(`it parses the template file relative to its base path
		and executes it with the given data`, func(t *testing.T) {
		writer := &strings.Builder{}
		err := loader.Load(writer, nil, "../../templates/sign-in.html")
		tests.AssertNoError(t, err)
	})

	t.Run(`it parses multiple template files relative to its base path
		and executes them with the given data`, func(t *testing.T) {
		writer := &strings.Builder{}
		err := loader.Load(writer, nil, "../../templates/app.html", "../../templates/sidebar.html")
		tests.AssertNoError(t, err)
	})

	t.Run("when it cannot load a template, it returns an error", func(t *testing.T) {
		writer := &strings.Builder{}
		err := loader.Load(writer, nil, "./unknown-template.html")
		tests.AssertError(t, err)
	})
}
