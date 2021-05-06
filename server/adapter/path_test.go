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
	"testing"

	"github.com/hyzual/mike-sierra-sierra/server/adapter"
)

func TestBasePathJoiner(t *testing.T) {
	basePath := "/path/to/app/assets"
	joiner := adapter.NewBasePathJoiner(basePath)

	t.Run("it joins relative paths to baseDir to form absolute paths", func(t *testing.T) {
		assertPathEquals(t, joiner.Join("./style.css"), "/path/to/app/assets/style.css")
	})

	t.Run("it joins nested paths to baseDir", func(t *testing.T) {
		assertPathEquals(t, joiner.Join("./sub/dir/style.css"), "/path/to/app/assets/sub/dir/style.css")
	})

	t.Run("it does not allow ascending up its base path", func(t *testing.T) {
		assertPathEquals(t, joiner.Join("../style.css"), "/path/to/app/assets")
	})

	t.Run("it does not allow ascending with a subdir before dot dot", func(t *testing.T) {
		assertPathEquals(t, joiner.Join("./sub/dir/../../../style.css"), "/path/to/app/assets")
	})
}

func assertPathEquals(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("joined path %s does not equal %s", got, want)
	}
}
