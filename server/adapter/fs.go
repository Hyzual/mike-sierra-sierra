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

package adapter

import (
	"io/fs"
	"os"

	"github.com/hyzual/mike-sierra-sierra/server/domain/music"
)

// NewOSFileSystem creates a new OSFileSystem
func NewOSFileSystem(dirFS fs.FS, pathJoiner PathJoiner) music.MusicLibraryFileSystem {
	return &baseOSFileSystem{dirFS, pathJoiner}
}

type baseOSFileSystem struct {
	dirFS      fs.FS
	pathJoiner PathJoiner
}

func (b *baseOSFileSystem) Open(name string) (fs.File, error) {
	return b.dirFS.Open(name)
}

func (b *baseOSFileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	path := b.pathJoiner.Join(name)
	return os.ReadDir(path)
}
