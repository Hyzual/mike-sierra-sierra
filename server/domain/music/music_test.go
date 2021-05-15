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

package music_test

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/hyzual/mike-sierra-sierra/server/domain/music"
	"github.com/hyzual/mike-sierra-sierra/tests"
)

func TestBaseMusicLibraryExplorer(t *testing.T) {
	t.Run("it returns an error when a given folder in the music library cannot be read", func(t *testing.T) {
		explorer := newExplorerWithError()
		_, _, err := explorer.ListContents("Sub Folder")
		tests.AssertError(t, err)
	})

	t.Run(`it sorts the contents of the given folder in the music library
		into songs and folders and returns them`, func(t *testing.T) {
		testFS := fstest.MapFS{
			"Sub Folder/against/sleep.mp3":       {},
			"Sub Folder/increase/furniture.flac": {},
			"Sub Folder/summer.mp3":              {},
			"Sub Folder/porch.ogg":               {},
		}
		explorer := music.NewMusicLibraryExplorer(testFS)
		folders, songs, err := explorer.ListContents("Sub Folder")

		tests.AssertNoError(t, err)
		assertSubFoldersContain(t, folders, music.SubFolder{Name: "against", Path: "Sub Folder/against"})
		assertSubFoldersContain(t, folders, music.SubFolder{Name: "increase", Path: "Sub Folder/increase"})

		assertSongsContain(t, songs, music.Song{Title: "summer.mp3", URI: "/music/Sub Folder/summer.mp3"})
		assertSongsContain(t, songs, music.Song{Title: "porch.ogg", URI: "/music/Sub Folder/porch.ogg"})
	})

	t.Run(`given "." as a path, it sorts the contents of the root music library folder
		into songs and folders and returns them`, func(t *testing.T) {
		testFS := fstest.MapFS{
			"ability/speech.mp3": {},
			"usual/compare.flac": {},
			"food shot.mp3":      {},
			"sit.flac":           {},
		}
		explorer := music.NewMusicLibraryExplorer(testFS)
		folders, songs, err := explorer.ListContents(".")

		tests.AssertNoError(t, err)
		assertSubFoldersContain(t, folders, music.SubFolder{Name: "ability", Path: "ability"})
		assertSubFoldersContain(t, folders, music.SubFolder{Name: "usual", Path: "usual"})

		assertSongsContain(t, songs, music.Song{Title: "food shot.mp3", URI: "/music/food shot.mp3"})
		assertSongsContain(t, songs, music.Song{Title: "sit.flac", URI: "/music/sit.flac"})
	})
}

func newExplorerWithError() music.MusicLibraryExplorer {
	return music.NewMusicLibraryExplorer(fstest.MapFS{
		".": {Mode: fs.FileMode(0000)},
	})
}

func assertSubFoldersContain(t *testing.T, slice []music.SubFolder, want music.SubFolder) {
	t.Helper()
	for _, got := range slice {
		if got.Name == want.Name {
			assertSubFolderEquals(t, got, want)
			return
		}
	}
	t.Errorf("could not find wanted sub-folder %v in slice %v", want, slice)
}

func assertSubFolderEquals(t *testing.T, got music.SubFolder, want music.SubFolder) {
	t.Helper()
	if got.Name != want.Name {
		t.Errorf("sub-folder name %s does not equal %s", got.Name, want.Name)
	}
	if got.Path != want.Path {
		t.Errorf("sub-folder path %s does not equal %s", got.Path, want.Path)
	}
}

func assertSongsContain(t *testing.T, slice []music.Song, want music.Song) {
	t.Helper()
	for _, got := range slice {
		if got.Title == want.Title {
			assertSongEquals(t, got, want)
			return
		}
	}
	t.Errorf("could not find wanted song %v in slice %v", want, slice)
}

func assertSongEquals(t *testing.T, got music.Song, want music.Song) {
	t.Helper()
	if got.Title != want.Title {
		t.Errorf("song title %s does not equal %s", got.Title, want.Title)
	}
	if got.URI != want.URI {
		t.Errorf("song URI %s does not equal %s", got.URI, want.URI)
	}
}
