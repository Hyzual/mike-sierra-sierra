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

import { errAsync, okAsync } from "neverthrow";
import * as rest_querier from "../api/rest-querier";
import type { Folder } from "../types";
import "./FolderDetails";
import { FolderDetails } from "./FolderDetails";

describe("FolderDetails", () => {
    const tag_name = "mss-folder-details";

    afterEach(() => {
        document.body.innerHTML = "";
    });

    it(`renders a loading state while loading the root folder`, async () => {
        jest.spyOn(rest_querier, "getFolder").mockReturnValue(
            okAsync({ folders: [], songs: [] })
        );
        const element = new FolderDetails();
        document.body.append(element);

        await element.updateComplete;
        expect(element.shadowRoot?.innerHTML).toContain("Loading ...");
    });

    it(`renders a list of folders once the root folder is loaded`, async () => {
        const async_result = okAsync<Folder, Error>({
            folders: [
                { name: "last", uri: "last" },
                { name: "direction", uri: "direction" },
            ],
            songs: [],
        });
        jest.spyOn(rest_querier, "getFolder").mockReturnValue(async_result);
        const element = new FolderDetails();
        document.body.append(element);

        await element.updateComplete;
        await async_result;
        expect(element.shadowRoot?.innerHTML).toContain("folders-list");
    });

    it(`given a folder_path property, it loads the folder at this path
        and renders a list of folders`, async () => {
        const async_result = okAsync<Folder, Error>({
            folders: [
                { name: "liquid", uri: "live/liquid" },
                { name: "wooden", uri: "live/wooden" },
            ],
            songs: [],
        });
        jest.spyOn(rest_querier, "getFolder").mockReturnValue(async_result);
        const element = new FolderDetails();
        element.folder_path = "live";
        document.body.append(element);

        await element.updateComplete;
        await async_result;
        expect(element.shadowRoot?.innerHTML).toContain("folders-list");
    });

    it(`when there is an error, it renders an error state`, async () => {
        const async_result = errAsync<Folder, Error>(
            new Error("Could not decode JSON")
        );
        jest.spyOn(rest_querier, "getFolder").mockReturnValue(async_result);
        const element = new FolderDetails();
        document.body.append(element);

        await element.updateComplete;
        await async_result;
        expect(element.shadowRoot?.innerHTML).toContain("An error occurred");
    });
});
