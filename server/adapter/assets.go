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

package adapter

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
)

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
	fs            fs.FS
	assetsBaseURI string
}

// NewAssetsResolver creates a new AssetsResolver
func NewAssetsResolver(fs fs.FS, assetsBaseURI string) AssetsResolver {
	return &baseAssetsResolver{fs, assetsBaseURI}
}

// GetAssetURI returns the asset's URI from its baseName.
// It reads the "manifest.json" file found at the root of baseAssetsResolver's fs and searches for baseName.
// It then joins the corresponding "hashed file name" (read from the manifest.json file)
// to its assetsBaseURI and returns the joined URI.
func (b *baseAssetsResolver) GetAssetURI(baseName string) (string, error) {
	const manifestPath string = "manifest.json"
	manifestFile, err := b.fs.Open(manifestPath)
	if err != nil {
		return "", fmt.Errorf("Could not read the manifest.json file: %w", err)
	}
	defer manifestFile.Close()

	var manifestContents assetsManifest
	err = json.NewDecoder(manifestFile).Decode(&manifestContents)
	if err != nil {
		return "", fmt.Errorf("Could not decode the manifest.json file: %w", err)
	}

	hashedFileName, ok := manifestContents[baseName]
	if !ok {
		return "", fmt.Errorf("Could not find %s in the manifest.json file", baseName)
	}
	joinedURI := path.Join(b.assetsBaseURI, hashedFileName)

	return joinedURI, nil
}

type assetsManifest = map[string]string
