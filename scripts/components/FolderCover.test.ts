/*
 *   Copyright (C) 2021  Joris MASSON
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

import { router } from "../router";
import { FolderCover } from "./FolderCover";

jest.mock("../../images/assets/no-cover.svg", () => {
    return "<svg></svg>";
});

const noop = (): void => {
    //Do nothing
};

describe(`FolderCover`, () => {
    let element: FolderCover;

    beforeEach(async () => {
        element = new FolderCover();
        element.setAttribute("folder_uri", "manufacturing/gently");
        element.setAttribute("folder_title", "Gently");
        router.on("/folders/:path", noop);

        document.body.append(element);
        await element.updateComplete;
    });

    afterEach(() => {
        element.remove();
    });

    it(`When I click on the link, it will navigate using the frontend router`, () => {
        const anchor = element.shadowRoot?.querySelector("a");
        if (!anchor) {
            throw new Error(
                "Could not find <a> tag in the element's shadowRoot"
            );
        }
        const navigate = jest.spyOn(router, "navigate");
        anchor.dispatchEvent(new Event("click"));
        expect(navigate).toHaveBeenCalledWith(
            "/folders/" + encodeURIComponent("manufacturing/gently")
        );
    });
});
