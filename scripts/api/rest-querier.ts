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

import { ono } from "@jsdevtools/ono";
import { ok, err, ResultAsync } from "neverthrow";
import { Folder, TopLevelFolders } from "../types";
import { NetworkError } from "./NetworkError";

const wrapError = (e: unknown): Error =>
    e instanceof Error ? e : ono("Unknown error");

export const getTopFolders = (): ResultAsync<
    TopLevelFolders,
    Error | NetworkError
> =>
    getAPI("/api/folders").andThen((response) =>
        ResultAsync.fromPromise(response.json(), wrapError).mapErr((error) =>
            ono(error, "Could not decode JSON into top-level folders")
        )
    );

export const getFolder = (
    folder_id: number
): ResultAsync<Folder, Error | NetworkError> =>
    getAPI(`/api/folders/${folder_id}`).andThen((response) =>
        ResultAsync.fromPromise(response.json(), wrapError).mapErr((error) =>
            ono(error, "Could not decode JSON into Folder")
        )
    );

function getAPI(uri: string): ResultAsync<Response, Error | NetworkError> {
    return ResultAsync.fromPromise(
        fetch(encodeURI(uri), { method: "GET" }),
        wrapError
    ).andThen((response) => {
        if (!response.ok) {
            return err(
                new NetworkError(
                    "GET",
                    uri,
                    response.status,
                    response.statusText
                )
            );
        }
        return ok(response);
    });
}
