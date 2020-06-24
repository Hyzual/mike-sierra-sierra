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
import { initSidebarLinks } from "./router";
import { init as initFontAwesome } from "./fontawesome";

const APP_MOUNT_POINT = "#app-mount-point";

document.addEventListener("DOMContentLoaded", () => {
    initFontAwesome();
    let error = initSidebarLinks(document, window);
    if (error) {
        throw error;
    }
    error = init(document);
    if (error) {
        throw error;
    }
    printFolder();
});

//TODO: UT
function init(doc: Document): Error | null {
    const mount_point = doc.querySelector(APP_MOUNT_POINT);
    if (!mount_point) {
        return new Error(
            `Could not get the app mount point at id ${APP_MOUNT_POINT}`
        );
    }

    return null;
}

async function printFolder(): Promise<void> {
    const res = await getFolder(0);
    if (res.isErr()) {
        // eslint-disable-next-line no-console
        console.log(res.error);
        return;
    }
    // eslint-disable-next-line no-console
    console.log(res.value);
}
