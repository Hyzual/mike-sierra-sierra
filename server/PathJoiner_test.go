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
	"testing"

	"github.com/hyzual/mike-sierra-sierra/server"
)

func TestRootPathJoiner(t *testing.T) {
	rootDir := "/path/to/app"
	joiner := server.NewRootPathJoiner(rootDir)

	t.Run("it joins relative paths to baseDir to form absolute paths", func(t *testing.T) {
		assertPathEquals(t, joiner.Join("./templates/login.html"), "/path/to/app/templates/login.html")
	})

	t.Run("it joins relative paths that ascend the hierarchy", func(t *testing.T) {
		assertPathEquals(t, joiner.Join("../assets/style.css"), "/path/to/assets/style.css")
	})
}

func assertPathEquals(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("joined path %s does not equal %s", got, want)
	}
}
