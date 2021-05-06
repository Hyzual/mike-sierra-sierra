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

package adapter_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"testing/fstest"

	"github.com/hyzual/mike-sierra-sierra/server/adapter"
	"github.com/hyzual/mike-sierra-sierra/tests"
)

func TestAssetsResolver(t *testing.T) {
	t.Run("when there is no manifest.json file in the assets directory, it will return an error", func(t *testing.T) {
		resolver := newResolverWithNoManifest(t)

		_, err := resolver.GetAssetURI("style.css")
		tests.AssertError(t, err)
	})

	t.Run("when the manifest.json file is not JSON-encoded, it will return an error", func(t *testing.T) {
		resolver := newResolverWithBadlyEncodedManifest(t)

		_, err := resolver.GetAssetURI("style.css")
		tests.AssertError(t, err)
	})

	t.Run("given a baseName not found in the manifest.json file, it will return an error", func(t *testing.T) {
		resolver := newResolverWithValidManifest(t)

		_, err := resolver.GetAssetURI("unknown/file.js")
		tests.AssertError(t, err)
	})

	t.Run(`given a baseName,
		it will read the hashed file name from the manifest.json file,
		and return it joined to its baseURI`, func(t *testing.T) {
		resolver := newResolverWithValidManifest(t)

		got, err := resolver.GetAssetURI("style.css")

		tests.AssertNoError(t, err)
		assertHashedNameEquals(t, got, "/assets/style.chunkhash.css")
	})

	t.Run("when baseName contains a slash, it will return the joined hashed file name without error", func(t *testing.T) {
		resolver := newResolverWithValidManifest(t)

		got, err := resolver.GetAssetURI("subdirectory/file.js")
		tests.AssertNoError(t, err)
		assertHashedNameEquals(t, got, "/assets/subdirectory/file.chunkhash.js")
	})
}

func newResolverWithNoManifest(t *testing.T) adapter.AssetsResolver {
	t.Helper()

	testFS := fstest.MapFS{
		"not-manifest.json": {},
	}
	// No manifest.json file

	return adapter.NewAssetsResolver(testFS, "/assets")
}

func newResolverWithBadlyEncodedManifest(t *testing.T) adapter.AssetsResolver {
	t.Helper()

	testFS := fstest.MapFS{
		"manifest.json": {Data: []byte{}},
	}
	// manifest is empty and does not contain JSON

	return adapter.NewAssetsResolver(testFS, "/assets")
}

func newResolverWithValidManifest(t *testing.T) adapter.AssetsResolver {
	manifestBuffer := new(bytes.Buffer)
	manifestContent := make(map[string]string)
	manifestContent["style.css"] = "style.chunkhash.css"
	manifestContent["subdirectory/file.js"] = "subdirectory/file.chunkhash.js"
	err := json.NewEncoder(manifestBuffer).Encode(manifestContent)
	if err != nil {
		t.Fatalf("Could not setup test manifest file, '%v'", err)
	}

	testFS := fstest.MapFS{
		"manifest.json": {Data: manifestBuffer.Bytes()},
	}
	return adapter.NewAssetsResolver(testFS, "/assets")
}

func assertHashedNameEquals(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("hashed file name %s does not equal %s", got, want)
	}
}
