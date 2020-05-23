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

/*
package walker groups together all code needed to walk the music library
and write info to the database for the webserver to read.
*/
package walker

import (
	"errors"

	"github.com/blang/vfs"
)

// MusicDirectoryWalker walks the music library and searches for music files. It reads them
// and writes information to the database for the webserver to read.
type MusicDirectoryWalker struct {
	fs    vfs.Filesystem
	store MusicLibraryStore
}

// MusicLibraryStore handles database operations related to the music library
// such as Folders, Albums, Songs.
type MusicLibraryStore interface {
	SaveFolder(info *folderInfo) error
	SaveSong(info *songInfo) error
}

type folderInfo struct {
	Base string
}

type songInfo struct {
	FolderID uint
	Title    string
}

// Walk scans the given path for music files. It reads them and
// writes information to the database
func (m *MusicDirectoryWalker) Walk(path string) error {
	return errors.New("Not implemented")
	//TODO: validate path ? Clean it ? or make sure it is absolute ?

	// info, err := os.Lstat(root)
	// if err != nil {
	// 	err = walkFn(root, nil, err)
	// } else {
	// 	err = walk(root, info, walkFn)
	// }
	// if err == SkipDir {
	// 	return nil
	// }
	// return err
}

// walk recursively descends path, calling walkFn.
// func walk(path string, info os.FileInfo, walkFn WalkFunc) error {
// 	if !info.IsDir() {
// 		return walkFn(path, info, nil)
// 	}

// 	names, err := readDirNames(path)
// 	err1 := walkFn(path, info, err)
// 	// If err != nil, walk can't walk into this directory.
// 	// err1 != nil means walkFn want walk to skip this directory or stop walking.
// 	// Therefore, if one of err and err1 isn't nil, walk will return.
// 	if err != nil || err1 != nil {
// 		// The caller's behavior is controlled by the return value, which is decided
// 		// by walkFn. walkFn may ignore err and return nil.
// 		// If walkFn returns SkipDir, it will be handled by the caller.
// 		// So walk should return whatever walkFn returns.
// 		return err1
// 	}

// 	for _, name := range names {
// 		filename := Join(path, name)
// 		fileInfo, err := lstat(filename)
// 		if err != nil {
// 			if err := walkFn(filename, fileInfo, err); err != nil && err != SkipDir {
// 				return err
// 			}
// 		} else {
// 			err = walk(filename, fileInfo, walkFn)
// 			if err != nil {
// 				if !fileInfo.IsDir() || err != SkipDir {
// 					return err
// 				}
// 			}
// 		}
// 	}
// 	return nil
// }

// func musicScanner(path string, info os.FileInfo, err error) error {

// }
