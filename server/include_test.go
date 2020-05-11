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

package server_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/blang/vfs"
	"github.com/blang/vfs/memfs"
	"github.com/hyzual/mike-sierra-sierra/server"
	"github.com/hyzual/mike-sierra-sierra/tests"
)

func TestBasePathJoiner(t *testing.T) {
	basePath := "/path/to/app/assets"
	joiner := server.NewBasePathJoiner(basePath)

	t.Run("it joins relative paths to baseDir to form absolute paths", func(t *testing.T) {
		assertPathEquals(t, joiner.Join("./style.css"), "/path/to/app/assets/style.css")
	})

	t.Run("it joins nested paths to baseDir", func(t *testing.T) {
		assertPathEquals(t, joiner.Join("./sub/dir/style.css"), "/path/to/app/assets/sub/dir/style.css")
	})

	t.Run("it does not allow ascending up its base path", func(t *testing.T) {
		assertPathEquals(t, joiner.Join("../style.css"), "/path/to/app/assets")
	})

	t.Run("it does not allow ascending with a subdir before dot dot", func(t *testing.T) {
		assertPathEquals(t, joiner.Join("./sub/dir/../../../style.css"), "/path/to/app/assets")
	})
}

func TestTemplateLoader(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not get the current working directory, '%v'", err)
	}
	loader := server.NewTemplateLoader(basePath)

	t.Run("it parses the template files relative to its base path", func(t *testing.T) {
		_, err := loader.Load("../templates/login.html")
		tests.AssertNoError(t, err)
	})

	t.Run("when it cannot load a template, it returns an error", func(t *testing.T) {
		_, err := loader.Load("./unknown-template.html")
		tests.AssertError(t, err)
	})
}

func TestAssetsResolver(t *testing.T) {

	t.Run("when there is no manifest.json file in the assets directory, it will return an error", func(t *testing.T) {
		resolver := newResolverWithNoManifest(t)

		_, err := resolver.GetHashedName("style.css")
		tests.AssertError(t, err)
	})

	t.Run("when the manifest.json file is not JSON-encoded, it will return an error", func(t *testing.T) {
		resolver := newResolverWithBadlyEncodedManifest(t)

		_, err := resolver.GetHashedName("style.css")
		tests.AssertError(t, err)
	})

	t.Run("given a baseName not found in the manifest.json file, it will return an error", func(t *testing.T) {
		resolver := newResolverWithValidManifest(t)

		_, err := resolver.GetHashedName("unknown/file.js")
		tests.AssertError(t, err)
	})

	t.Run("given a baseName, it will return the hashed file name from the manifest.json file", func(t *testing.T) {
		resolver := newResolverWithValidManifest(t)

		got, err := resolver.GetHashedName("style.css")

		tests.AssertNoError(t, err)
		assertHashedNameEquals(t, got, "style.chunkhash.css")
	})

	t.Run("when baseName contains a slash, it will return the hashed file name without error", func(t *testing.T) {
		resolver := newResolverWithValidManifest(t)

		got, err := resolver.GetHashedName("subdirectory/file.js")
		tests.AssertNoError(t, err)
		assertHashedNameEquals(t, got, "subdirectory/file.chunkhash.js")
	})
}

func newResolverWithNoManifest(t *testing.T) server.AssetsResolver {
	t.Helper()

	testFS := memfs.Create()
	testFS.Mkdir("/assets", 0755)
	// No manifest.json file

	return server.NewAssetsResolver(testFS, "/assets")
}

func buildManifestFile(t *testing.T) (vfs.Filesystem, vfs.File, func()) {
	t.Helper()

	testFS := memfs.Create()
	testFS.Mkdir("/assets", 0755)
	manifest, err := testFS.OpenFile("/assets/manifest.json", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		t.Fatalf("Could not setup test manifest file, '%v'", err)
	}

	return testFS, manifest, func() {
		manifest.Close()
	}
}

func newResolverWithBadlyEncodedManifest(t *testing.T) server.AssetsResolver {
	t.Helper()

	testFS, _, closeManifestFile := buildManifestFile(t)
	defer closeManifestFile()
	// manifest is empty and does not contain JSON

	return server.NewAssetsResolver(testFS, "/assets")
}

func newResolverWithValidManifest(t *testing.T) server.AssetsResolver {
	testFS, manifest, closeManifestFile := buildManifestFile(t)
	defer closeManifestFile()
	var manifestContent = make(map[string]string)
	manifestContent["style.css"] = "style.chunkhash.css"
	manifestContent["subdirectory/file.js"] = "subdirectory/file.chunkhash.js"
	err := json.NewEncoder(manifest).Encode(manifestContent)
	if err != nil {
		t.Fatalf("Could not setup test manifest file, '%v'", err)
	}
	return server.NewAssetsResolver(testFS, "/assets")
}

func assertPathEquals(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("joined path %s does not equal %s", got, want)
	}
}

func assertHashedNameEquals(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("hashed file name %s does not equal %s", got, want)
	}
}
