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

func TestAssetsJoin(t *testing.T) {
	basePath := "/path/to/app/assets"
	includer := server.NewAssetsIncluder(basePath)

	t.Run("it joins relative paths to baseDir to form absolute paths", func(t *testing.T) {
		assertPathEquals(t, includer.Join("./style.css"), "/path/to/app/assets/style.css")
	})

	//TODO: it should not allow ascending
	t.Run("it joins relative paths that ascend the hierarchy", func(t *testing.T) {
		assertPathEquals(t, includer.Join("../style.css"), "/path/to/app/style.css")
	})
}

func TestMusicJoin(t *testing.T) {
	basePath := "/path/to/music"
	loader := server.NewMusicLoader(basePath)

	t.Run("it joins relative paths to baseDir to form absolute paths", func(t *testing.T) {
		assertPathEquals(t, loader.Join("./song.mp3"), "/path/to/music/song.mp3")
	})

	//TODO: it should not allow ascending
	t.Run("it joins relative paths that ascend the hierarchy", func(t *testing.T) {
		assertPathEquals(t, loader.Join("../song.mp3"), "/path/to/song.mp3")
	})
}

func assertPathEquals(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("joined path %s does not equal %s", got, want)
	}
}
