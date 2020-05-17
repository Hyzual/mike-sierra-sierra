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
	"encoding/json"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/blang/vfs"
	"github.com/pkg/errors"
)

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

// TemplateExecutor resolves the given relative file path from the templates/ folder,
// parses the template from the file and executes it on the given writer.
// It returns an error if the template can't be found or if the execution fails.
type TemplateExecutor interface {
	Load(writer io.Writer, templatePath string, data interface{}) error
}

// templateBaseExecutor implements TemplateExecutor for production code
type templateBaseExecutor struct {
	basePath string // absolute path to the /templates directory
}

// NewTemplateExecutor creates a new TemplateExecutor
func NewTemplateExecutor(basePath string) TemplateExecutor {
	return &templateBaseExecutor{basePath}
}

// Load loads the template at templatePath (relative path from the templates/ folder),
// parses the template from the file and executes it on the given writer.
// It returns an error if the template can't be found or if the execution fails.
func (t *templateBaseExecutor) Load(writer io.Writer, templatePath string, data interface{}) error {
	cleanedPath := path.Join(t.basePath, filepath.Clean(templatePath))
	tmpl, err := template.ParseFiles(cleanedPath)
	if err != nil {
		return errors.Wrapf(err, "could not load the template %s", templatePath)
	}
	err = tmpl.Execute(writer, data)
	if err != nil {
		return errors.Wrapf(err, "could not execute the template %s", templatePath)
	}
	return nil
}

// AssetsResolver reads the manifest.json file in the /assets directory
// It is used by templates to resolve asset URIs with chunkhashes in their names.
// For example, it translates "style.css" to "/assets/style-caf5894036274013394c.css".
// It returns an error if it can't find the manifest file or can't decode its contents
// from JSON or can't find the given baseName in the manifest contents
type AssetsResolver interface {
	GetAssetURI(baseName string) (string, error)
}

// baseAssetsResoler inmplements AssetsResolver
type baseAssetsResolver struct {
	fs             vfs.Filesystem
	assetsBasePath string
	assetsBaseURI  string
}

// NewAssetsResolver creates a new AssetsResolver
func NewAssetsResolver(fs vfs.Filesystem, assetsBasePath, assetsBaseURI string) AssetsResolver {
	return &baseAssetsResolver{fs, assetsBasePath, assetsBaseURI}
}

// GetAssetURI returns the asset's URI from its baseName.
// It reads the "manifest.json" file found at assetsBasePath and searches for baseName.
// It then joins the corresponding "hashed file name" (read from the manifest.json file)
// to its assetsBaseURI and returns the joined URI.
func (b *baseAssetsResolver) GetAssetURI(baseName string) (string, error) {
	manifestPath := path.Join(b.assetsBasePath, "./manifest.json")
	manifestFile, err := b.fs.OpenFile(manifestPath, os.O_RDONLY, 0)
	if err != nil {
		return "", errors.Wrapf(err, "Could not read the manifest.json file in this folder: %s. Did you run 'npm run build' ?", b.assetsBasePath)
	}
	defer manifestFile.Close()

	var manifestContents assetsManifest
	err = json.NewDecoder(manifestFile).Decode(&manifestContents)
	if err != nil {
		return "", errors.Wrap(err, "Could not decode the manifest.json file")
	}

	hashedFileName, ok := manifestContents[baseName]
	if !ok {
		return "", errors.Errorf("Could not find %s in the manifest.json file", baseName)
	}
	joinedURI := path.Join(b.assetsBaseURI, hashedFileName)

	return joinedURI, nil
}

type assetsManifest = map[string]string
