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

package server

import "path/filepath"

//PathJoiner joins filesystem paths to its base path
type PathJoiner interface {
	Join(string) string
}

//NewRootPathJoiner returns a new RootPathJoiner
func NewRootPathJoiner(rootDirectoryPath string) PathJoiner {
	return &RootPathJoiner{rootDirectoryPath}
}

//RootPathJoiner implements PathJoiner with its base path set to the project's root
type RootPathJoiner struct {
	rootDirectoryPath string //Project root directory from which to load templates and assets as a relative path
}

//Join the given path to the project's root directory
func (r *RootPathJoiner) Join(path string) string {
	return filepath.Join(r.rootDirectoryPath, path)
}
