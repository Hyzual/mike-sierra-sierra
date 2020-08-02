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

import { SideBarLink } from "./SideBarLink";
import { router } from "../router";

describe(`SideBarLink`, () => {
    let link: SideBarLink;

    beforeEach(async () => {
        link = new SideBarLink();
        link.setAttribute("uri", "/folders");
        link.setAttribute("label", "Browse by Folder");

        document.body.append(link);
        await link.updateComplete;
    });

    afterEach(() => {
        link.remove();
    });

    it(`Given a uri and a label, it will render a sidebar link`, () => {
        expect(
            link.shadowRoot?.querySelector("a")?.outerHTML
        ).toMatchInlineSnapshot(
            `"<a class=\\"link\\" href=\\"/folders\\"><slot name=\\"icon\\"></slot>Browse by Folder<!----></a>"`
        );
    });

    it(`When I click on the link, it will navigate using the frontend router`, () => {
        const anchor = link.shadowRoot?.querySelector("a");
        if (!anchor) {
            throw new Error(
                "Could not find <a> tag in the element's shadowRoot"
            );
        }
        const navigate = jest.spyOn(router, "navigate");
        anchor.dispatchEvent(new Event("click"));
        expect(navigate).toHaveBeenCalledWith("/folders");
    });
});
