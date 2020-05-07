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

import (
	"html/template"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

// PathJoiner joins the given relative path to its base path
type PathJoiner interface {
	Join(relativePath string) string
}

// AssetsIncluder implements PathJoiner
type AssetsIncluder interface {
	PathJoiner
}

// AssetsBaseIncluder implements AssetsIncluder for production code
type AssetsBaseIncluder struct {
	basePath string // absolute path to the /assets directory
}

// NewAssetsIncluder creates a new AssetsIncluder
func NewAssetsIncluder(basePath string) AssetsIncluder {
	return &AssetsBaseIncluder{basePath}
}

// Join joins the given relative path to the its base path
func (a *AssetsBaseIncluder) Join(relativePath string) string {
	return path.Join(a.basePath, relativePath)
}

// MusicBaseLoader implements PathJoiner
type MusicBaseLoader struct {
	basePath string // absolute path to the /music directory
}

// NewMusicLoader creates a new PathJoiner
func NewMusicLoader(basePath string) PathJoiner {
	return &MusicBaseLoader{basePath}
}

// Join joins the given relative path to the its base path
func (m *MusicBaseLoader) Join(relativePath string) string {
	return path.Join(m.basePath, relativePath)
}

// TemplateLoader resolves the given relative file path from the templates/ folder
// and returns the template parsed from the file or an error
type TemplateLoader interface {
	Load(templatePath string) (*template.Template, error)
}

// TemplateBaseLoader implements TemplateLoader for production code
type TemplateBaseLoader struct {
	basePath string // absolute path to the /templates directory
}

// NewTemplateLoader creates a new TemplateLoader
func NewTemplateLoader(basePath string) TemplateLoader {
	return &TemplateBaseLoader{basePath}
}

// Load loads the template at templatePath (relative path from the templates/ folder)
// and returns it. It returns an error if any (for example no file exists at templatePath)
func (t *TemplateBaseLoader) Load(templatePath string) (*template.Template, error) {
	cleanedPath := path.Join(t.basePath, filepath.Clean(templatePath))
	tmpl, err := template.ParseFiles(cleanedPath)

	if err != nil {
		return nil, errors.Wrap(err, "could not load template")
	}
	return tmpl, nil
}
