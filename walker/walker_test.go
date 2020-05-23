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

package walker

import (
	"errors"
	"os"
	"testing"

	"github.com/hyzual/mike-sierra-sierra/tests"

	"github.com/blang/vfs/memfs"
)

func TestWalker(t *testing.T) {
	//TODO: set up the happy-path music folder
	testFS := memfs.Create()
	err := testFS.Mkdir("/music/", 0755)
	if err != nil {
		t.Fatalf("Could not setup test music folder: %v", err)
	}
	//TODO: use a fixture file ?
	songFile, err := testFS.OpenFile("/music/song.mp3", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	songFile.Write("dude")

	store := &stubMusicLibraryStore{}
	walker := &MusicDirectoryWalker{testFS, store}

	err = walker.Walk("/music")

	tests.AssertNoError(t, err)
}

type stubMusicLibraryStore struct{}

func (s *stubMusicLibraryStore) SaveFolder(info *folderInfo) error {
	return errors.New("This method should not have been called in the test")
}

func (s *stubMusicLibraryStore) SaveSong(info *songInfo) error {
	return errors.New("This method should not have been called in the test")
}
