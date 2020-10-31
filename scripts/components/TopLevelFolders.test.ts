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

import { LitElement } from "lit-element";
import { errAsync, okAsync } from "neverthrow";
import { TopLevelFolders } from "scripts/types";
import * as rest_querier from "../api/rest-querier";
import "./TopLevelFolders";

describe("TopLevelFolders", () => {
    const tag_name = "mss-top-level-folders";

    afterEach(() => {
        document.body.innerHTML = "";
    });

    it(`renders a loading state while loading the top-level folders`, async () => {
        jest.spyOn(rest_querier, "getTopFolders").mockReturnValue(
            okAsync({ folders: [] })
        );
        const element = document.createElement(tag_name) as LitElement;
        document.body.appendChild(element);

        await element.updateComplete;
        expect(element.shadowRoot?.innerHTML).toContain("Loading ...");
    });

    it(`renders a list of folders once the top-level folders are loaded`, async () => {
        const async_result = okAsync<TopLevelFolders, Error>({
            folders: [
                { id: 1, name: "last" },
                { id: 2, name: "direction" },
            ],
        });
        jest.spyOn(rest_querier, "getTopFolders").mockReturnValue(async_result);
        const element = document.createElement(tag_name) as LitElement;
        document.body.appendChild(element);

        await element.updateComplete;
        await async_result;
        expect(element.shadowRoot?.innerHTML).toContain("folders-list");
    });

    it(`when there is an error, it renders an error state`, async () => {
        const async_result = errAsync<TopLevelFolders, Error>(
            new Error("Could not decode JSON")
        );
        jest.spyOn(rest_querier, "getTopFolders").mockReturnValue(async_result);
        const element = document.createElement(tag_name) as LitElement;
        document.body.appendChild(element);

        await element.updateComplete;
        await async_result;
        expect(element.shadowRoot?.innerHTML).toContain("An error occurred");
    });
});
