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

import { initSidebarLinks } from "./router";

describe(`router`, () => {
    let doc: Document;
    beforeEach(() => {
        doc = document.implementation.createHTMLDocument();
    });

    describe(`initSidebarLinks()`, () => {
        describe(`Given a selector for the sidebar menu`, () => {
            it(`will return an error when the selector can't be found`, () => {
                const error = initSidebarLinks(doc, window);
                expect(error).not.toBeNull();
            });

            it(`will loop over its children, search for all [data-router-link]
                and if a tag that is not <a> is found, it will return an error`, () => {
                const sidebar_menu = doc.createElement("ul");
                sidebar_menu.id = "sidebar-menu-app-links";
                const li = doc.createElement("li");
                const not_anchor_tag = doc.createElement("span");
                not_anchor_tag.dataset.routerLink = "";
                li.append(not_anchor_tag);
                sidebar_menu.append(li);
                doc.body.append(sidebar_menu);

                const error = initSidebarLinks(doc, window);

                expect(error).not.toBeNull();
            });

            it(`will loop over its children, search for all a[data-router-link] tags
                and add a click listener on each.
                When the link is clicked, it will prevent default
                and it will push an URI in history`, () => {
                const sidebar_menu = doc.createElement("ul");
                sidebar_menu.id = "sidebar-menu-app-links";
                const li = doc.createElement("li");
                const anchor_tag = doc.createElement("a");
                anchor_tag.href = "/app/folders";
                anchor_tag.dataset.routerLink = "";
                li.append(anchor_tag);
                sidebar_menu.append(li);
                doc.body.append(sidebar_menu);

                const error = initSidebarLinks(doc, window);
                expect(error).toBeNull();

                const pushState = jest.spyOn(window.history, "pushState");
                anchor_tag.dispatchEvent(new Event("click"));
                expect(pushState).toHaveBeenCalledWith("", "", "/app/folders");
            });
        });
    });
});
