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

import {
    getFolder,
    getTopFolders,
    Folder,
    TopLevelFolders,
} from "./rest-querier";
import { ResultAsync } from "neverthrow";
import { NetworkError } from "./NetworkError";

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

    it.each([
        [
            "getFolder()",
            "Could not GET /api/folders/1",
            (): ResultAsync<unknown, NetworkError> => getFolder(1),
        ],
        [
            "getTopFolders()",
            "Could not GET /api/folders",
            (): ResultAsync<unknown, NetworkError> => getTopFolders(),
        ],
    ])(
        `when there is a network error, %s will return an error`,
        async (
            _test_name: string,
            expected_error_message: string,
            functionUnderTest
        ) => {
            globalFetch.mockImplementation(() =>
                Promise.resolve({
                    ok: false,
                    statusText: "Not found",
                })
            );

            const result = await functionUnderTest();
            if (!result.isErr()) {
                throw new Error("Expected an error but did not get one");
            }
            expect(result.error.message).toMatch(expected_error_message);
        }
    );

    it.each([
        [
            "getFolder()",
            "Could not decode JSON into Folder",
            (): ResultAsync<unknown, Error> => getFolder(1),
        ],
        [
            "getTopFolders",
            "Could not decode JSON into top-level folders",
            (): ResultAsync<unknown, Error> => getTopFolders(),
        ],
    ])(
        `when the JSON response cannot be decoded, %s will return an error`,
        async (
            _test_name: string,
            expected_error_message: string,
            functionUnderTest
        ) => {
            globalFetch.mockImplementation(() =>
                Promise.resolve({
                    ok: true,
                    json: () => Promise.reject(new Error("Error in JSON")),
                })
            );

            const result = await functionUnderTest();
            if (!result.isErr()) {
                throw new Error("Expected an error but did not get one");
            }
            expect(result.error.message).toMatch(expected_error_message);
        }
    );

    it(`getFolder() will return a Folder`, async () => {
        const expected_folder: Folder = {
            name: "Edith Chapman",
            items: [],
        };
        mockFetchSuccess(expected_folder);

        const result = await getFolder(0);
        if (!result.isOk()) {
            throw new Error("Did not expect an error but got one");
        }
        expect(result.value).toEqual(expected_folder);
    });

    it(`getTopFolders() will return the top-level Folders`, async () => {
        const expected_top_level_folders: TopLevelFolders = {
            folders: [
                { id: 1, name: "pharos" },
                { id: 2, name: "elution" },
            ],
        };
        mockFetchSuccess(expected_top_level_folders);

        const result = await getTopFolders();
        if (!result.isOk()) {
            throw new Error("Did not expect an error but got one");
        }
        expect(result.value).toEqual(expected_top_level_folders);
    });
});
