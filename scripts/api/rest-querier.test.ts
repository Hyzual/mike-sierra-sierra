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

import { getFolder } from "./rest-querier";
import { Folder } from "scripts/types";

describe(`rest-querier`, () => {
    let globalFetch: jest.SpyInstance;
    beforeEach(() => {
        window.fetch = globalFetch = jest.fn();
    });

    afterEach(() => {
        window.fetch = (): Promise<Response> => {
            throw new Error("Not supposed to happen");
        };
    });

    function mockFetchSuccess(return_json: unknown, headers = {}): void {
        globalFetch.mockImplementation(() =>
            Promise.resolve({
                headers,
                ok: true,
                json: () => Promise.resolve(return_json),
            })
        );
    }

    it(`when there is a network error, getFolder() will return an error`, async () => {
        globalFetch.mockImplementation(() =>
            Promise.resolve({
                ok: false,
                statusText: "Not found",
            })
        );

        const result = await getFolder("path");
        if (!result.isErr()) {
            throw new Error("Expected an error but did not get one");
        }
        expect(result.error.message).toMatch("Could not GET /api/folders/path");
    });

    it(`when the JSON response cannot be decoded, getFolder() will return an error`, async () => {
        globalFetch.mockImplementation(() =>
            Promise.resolve({
                ok: true,
                json: () => Promise.reject(new Error("Error in JSON")),
            })
        );

        const result = await getFolder("path");
        if (!result.isErr()) {
            throw new Error("Expected an error but did not get one");
        }
        expect(result.error.message).toMatch(
            "Could not decode JSON into Folder"
        );
    });

    it(`getFolder("path") will return a Folder`, async () => {
        const expected_folder: Folder = {
            folders: [],
            songs: [],
        };
        mockFetchSuccess(expected_folder);

        const result = await getFolder("path");
        if (!result.isOk()) {
            throw new Error("Did not expect an error but got one");
        }
        expect(result.value).toEqual(expected_folder);
    });

    it(`getFolder("") will return the top-level (root) music Folder`, async () => {
        const expected_folder: Folder = {
            folders: [
                { name: "pharos", uri: "pharos" },
                { name: "elution", uri: "elution" },
            ],
            songs: [{ title: "shall" }, { title: "asleep" }],
        };
        mockFetchSuccess(expected_folder);

        const result = await getFolder("");
        if (!result.isOk()) {
            throw new Error("Did not expect an error but got one");
        }
        expect(result.value).toEqual(expected_folder);
    });
});
