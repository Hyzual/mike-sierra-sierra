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

import { SidebarLink, navigate } from "./SidebarLink";
import type { SidebarLinkElement } from "./SidebarLink";
import { router } from "../../router";

const noop = (): void => {
    //Do nothing
};

describe(`SidebarLink`, () => {
    let host: SidebarLinkElement & HTMLElement;

    beforeEach(() => {
        host = {
            uri: "/folders",
            label: "Browe by folder",
        } as SidebarLinkElement & HTMLElement;

        router.on("/folders", noop);
    });

    it(`Given a uri and a label, it will render a sidebar link`, () => {
        const doc = document.implementation.createHTMLDocument();
        const text_node = doc.createTextNode("");

        if (SidebarLink.render === undefined) {
            throw new Error("Hybrids element should have a render function");
        }

        const update = SidebarLink.render(host);
        update(host, text_node);

        expect(text_node.textContent).toMatchInlineSnapshot(
            `"<a href=\\"/app/folders\\"><slot></slot><!---->Browse by Folder</a>"`
        );
    });

    it(`When I click on the link, it will navigate using the frontend router`, () => {
        const routerNavigate = jest.spyOn(router, "navigate");
        navigate(host, new Event("click"));
        expect(routerNavigate).toHaveBeenCalledWith("/folders");
    });
});
