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

import { getFolder } from "./rest-querier";

describe(`rest-querier`, () => {
    let fetchMock: jest.SpyInstance;
    beforeEach(() => {
        // fetch does not exist in jsdom, so we must create it
        // eslint-disable-next-line jest/prefer-spy-on
        fetchMock = window.fetch = jest.fn();
    });

    afterEach(() => {
        delete window.fetch;
    });

    describe(`getFolder()`, () => {
        it(`when there is a network error, it will return an error`, async () => {
            const expected_response = {
                ok: false,
                statusText: "Not found",
            };
            fetchMock.mockImplementation(() =>
                Promise.resolve(expected_response)
            );

            const result = await getFolder(0);
            if (!result.isErr()) {
                throw new Error("Expected an error but did not get one");
            }
            expect(result.error).toHaveProperty(
                "message",
                "Could not GET /api/folders/0"
            );
        });

        it(`when there is a JSON decoding error, it will return an error`, async () => {
            fetchMock.mockImplementation(() =>
                Promise.resolve({
                    ok: true,
                    json: () => Promise.reject(new Error("Error in JSON")),
                })
            );

            const result = await getFolder(0);
            if (!result.isErr()) {
                throw new Error("Expected an error but did not get one");
            }
            expect(result.error.message).toMatch("Could not decode JSON into Folder");
        });

        it(`will return a JSON Folder`, async () => {
            const expected_folder = {
                name: "Edith Chapman",
                items: []
            };
            mockFetchSuccess(expected_folder);

            const result = await getFolder(0);
            if (!result.isOk()) {
                throw new Error("Did not expect an error but got one");
            }
            expect(result.value).toEqual(expected_folder);
        });
    });

    function mockFetchSuccess(return_json: object, headers = {}): void {
        fetchMock.mockImplementation(() =>
            Promise.resolve({
                headers,
                ok: true,
                json: () => Promise.resolve(return_json),
            })
        );
    }
});
