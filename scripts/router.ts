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

import Navigo from "navigo";

const SIDEBAR_MENU_APP_LINKS = "#sidebar-menu-app-links";

export const router = new Navigo();

export function initSidebarLinks(doc: Document): Error | null {
    const sidebar_menu = doc.querySelector(SIDEBAR_MENU_APP_LINKS);
    if (!sidebar_menu) {
        return new Error(
            `Could not get the sidebar app menu at id ${SIDEBAR_MENU_APP_LINKS}`
        );
    }

    const sidebar_links = sidebar_menu.querySelectorAll("[data-router-link]");
    for (const link of sidebar_links) {
        if (!(link instanceof HTMLAnchorElement)) {
            return new Error(
                `[data-router-link] should only be applied to <a> tags, got ${link.tagName}`
            );
        }

        const location = link.getAttribute("href");
        if (location === null) {
            return new Error(
                `Could not get the href attribute of a data-router-link <a> tag`
            );
        }

        //TODO: mark the active link with css
        link.addEventListener("click", (event: Event) => {
            event.preventDefault();
            router.navigate(location);
        });
    }

    return null;
}
