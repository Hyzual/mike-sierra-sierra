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

package music

import (
	"fmt"
	"io/fs"
	"path"
	"strings"
)

const MusicPath = "/music" // It is a volume in the Docker image

// SubFolder represents a folder that is either top-level or a child of a top-level
// music folder. A SubFolder can have zero or many children SubFolders.
type SubFolder struct {
	Name string // Basename of the folder. For example "Dark Passion Play"
	Path string // Absolute path to the folder. For example "/music/Symphonic Metal/Nightwish/Dark Passion Play/"
}

// Song represents a music file. It is distinguished by media type (audio/mp3, audio/flac, etc.)
// Most of its fields mirror tags such as ID3 tags for MP3.
type Song struct {
	Title string // Title of the song
	URI   string // URI to play the song. For example "/music/Symphonic Metal/Nightwish/Dark Passion Play/7 Days to the Wolves.ogg"
}

// MusicLibraryExplorer allows to explore the contents of the music library folders. It needs a MusicLibraryFileSystem.
// It transforms Directories and Files into Folders and Songs structs.
type MusicLibraryExplorer interface {
	// ListContents reads a given folder in the music library and sorts children into folders
	// (directories) and songs. It ignores other files. Given "." as folderPath, it will return
	// the contents of the root music library folder.
	ListContents(folderPath string) ([]SubFolder, []Song, error)
}

// MusicLibraryFileSystem allows to read the music library's root folder and to open the music files.
// It does not allow to open parents of the music library's root folder.
type MusicLibraryFileSystem = fs.ReadDirFS

// baseMusicLibraryExplorer implements MusicLibraryExplorer
type baseMusicLibraryExplorer struct {
	filesystem MusicLibraryFileSystem
}

func (b *baseMusicLibraryExplorer) ListContents(folderPath string) ([]SubFolder, []Song, error) {
	var (
		subFolders []SubFolder
		songs      []Song
	)
	entries, err := b.filesystem.ReadDir(folderPath)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not read the %v directory: %w", folderPath, err)
	}
	for _, entry := range entries {
		filePath := path.Join(folderPath, entry.Name())
		if entry.IsDir() {
			subFolders = append(subFolders, SubFolder{Name: entry.Name(), Path: filePath})
		} else if isFileASong(entry) {
			songURI := path.Join(MusicPath, filePath)
			songs = append(songs, Song{Title: entry.Name(), URI: songURI})
		}
	}
	return subFolders, songs, nil
}

var supportedExtensions = [3]string{".mp3", ".flac", ".ogg"}

func isFileASong(entry fs.DirEntry) bool {
	for _, extension := range supportedExtensions {
		if strings.HasSuffix(entry.Name(), extension) {
			return true
		}
	}
	return false
}

// NewMusicLibraryExplorer creates a new MusicLibraryExplorer
func NewMusicLibraryExplorer(filesystem fs.ReadDirFS) MusicLibraryExplorer {
	return &baseMusicLibraryExplorer{filesystem}
}
