/*
 *   Copyright (C) 2021  Joris MASSON
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
	"io/fs"
	"os"
	"path"
	"testing"
	"testing/fstest"

	"github.com/hyzual/mike-sierra-sierra/server/adapter"
	"github.com/hyzual/mike-sierra-sierra/server/domain/music"
	"github.com/hyzual/mike-sierra-sierra/tests"
)

func TestOSFileSystem(t *testing.T) {
	tempDir := t.TempDir()
	pathToFolder := path.Join(tempDir, "foo/bar/folder")
	err := os.MkdirAll(pathToFolder, fs.ModePerm)
	if err != nil {
		t.Fatalf("Could not create temporary testing directory at %v", pathToFolder)
	}

	t.Run("it returns an error when given name does not exist", func(t *testing.T) {
		fs := newFSWithPathJoiner(tempDir)

		_, err := fs.ReadDir("unknown")
		tests.AssertError(t, err)
	})

	t.Run("it does not allow to read a parent directory to its root", func(t *testing.T) {
		fs := newFSWithPathJoiner(pathToFolder)

		entries, err := fs.ReadDir("..")
		tests.AssertNoError(t, err)
		for _, entry := range entries {
			if entry.Name() != "folder" {
				t.Errorf("Could access to parent of root directory %v", pathToFolder)
			}
		}
	})

	t.Run("it allow to read a sub-folder", func(t *testing.T) {
		fs := newFSWithPathJoiner(tempDir)

		entries, err := fs.ReadDir("foo")
		tests.AssertNoError(t, err)
		for _, entry := range entries {
			if entry.Name() != "bar" {
				t.Errorf("Could not list sub-folders of root directory %v", tempDir)
			}
		}
	})

	t.Run("it allows to read a descendant folder", func(t *testing.T) {
		fs := newFSWithPathJoiner(tempDir)

		entries, err := fs.ReadDir("foo/bar")
		tests.AssertNoError(t, err)
		for _, entry := range entries {
			if entry.Name() != "folder" {
				t.Errorf("Could not list descendant folder of root directory %v", tempDir)
			}
		}
	})

	t.Run("it allows to Open a file", func(t *testing.T) {
		fs := newFSWithMap(pathToFolder)

		_, err := fs.Open("file.mp3")
		tests.AssertNoError(t, err)
	})
}

func newFSWithPathJoiner(basepath string) music.MusicLibraryFileSystem {
	return adapter.NewOSFileSystem(fstest.MapFS{}, adapter.NewBasePathJoiner(basepath))
}

func newFSWithMap(basepath string) music.MusicLibraryFileSystem {
	return adapter.NewOSFileSystem(fstest.MapFS{"file.mp3": {}}, adapter.NewBasePathJoiner(basepath))
}
