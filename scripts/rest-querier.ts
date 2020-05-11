/*
 *   Copyright (c) 2020 Joris MASSON

 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.

 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.

 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import { ono } from "@jsdevtools/ono";

interface Song {
    title: string;
}

interface Folder {
    name: string;
    items: Song[];
}

interface RESTResult<T> {
    error: Error | null;
    result: T | null;
}

export async function getFolder(folder_id: number): Promise<RESTResult<Folder>> {
    const response = await fetch(`/api/folders/${folder_id}`, {
        method: "GET",
        headers: new Headers(),
    });
    if (!response.ok) {
        return {
            error: ono.error(`Could not GET /api/folders/${folder_id}`),
            result: null,
        };
    }
    try {
        const folder = await response.json();
        return { error: null, result: folder };
    } catch (err) {
        return {
            error: ono(err, "Could not decode JSON into Folder"),
            result: null,
        };
    }
}
