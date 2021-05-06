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

/*
Package adapter holds all code dedicated to side-effects and communication with
the outside world (databse, os, filesystem, REST API, etc.).
The domain MUST NEVER depend on adapters, always the reverse: the adapters depend
on the domain.
*/
package adapter

import "path"

// PathJoiner joins the given relative path to its base path
type PathJoiner interface {
	Join(relativePath string) string
}

// NewBasePathJoiner creates a new PathJoiner
func NewBasePathJoiner(basePath string) PathJoiner {
	return &basePathJoiner{basePath}
}

// basePathJoiner implements PathJoiner. It is given a base path and will
// Join all relative paths to it.
type basePathJoiner struct {
	basePath string // absolute path
}

// Join joins the given relative path to the basePath
func (b *basePathJoiner) Join(relativePath string) string {
	dir := path.Dir(relativePath)
	if dir == ".." {
		return b.basePath
	}
	return path.Join(b.basePath, relativePath)
}
