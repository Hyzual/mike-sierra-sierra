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

import type { PropertyDeclarations, TemplateResult } from "lit";
import { css, html, LitElement } from "lit";
import { until } from "lit-html/directives/until";
import { NetworkError } from "../../api/NetworkError";
import { getFolder } from "../../api/rest-querier";
import type { PlayQueueState } from "../music/PlayQueueState";

const renderErrorState = (error: Error | NetworkError): TemplateResult => {
    const error_template =
        error instanceof NetworkError
            ? html`<p class="error">
                  An error occurred: Code ${error.statusCode}:
                  ${error.statusText}
              </p>`
            : html`<p class="error">An error occurred: ${error.message}</p>`;
    return html`<div class="error-container">${error_template}</div>`;
};

const renderFolder = async (
    play_queue: PlayQueueState,
    folder_path: string
): Promise<TemplateResult> => {
    const result = await getFolder(folder_path);
    if (result.isErr()) {
        return renderErrorState(result.error);
    }
    return html`<mss-folders-list
            .folders=${result.value.folders}
        ></mss-folders-list>
        <mss-songs-list
            .songs=${result.value.songs}
            .play_queue=${play_queue}
        ></mss-songs-list>`;
};

export class FolderDetails extends LitElement {
    folder_path = "";
    play_queue!: PlayQueueState;

    static get properties(): PropertyDeclarations {
        return { folder_path: { type: String } };
    }

    static readonly styles = css`
        .error-container {
            display: grid;
            height: 100%;
            align-items: center;
            background-color: var(--error-color);
        }

        .error {
            margin: 0 0 0 16px;
        }
    `;

    render(): TemplateResult {
        return html`${until(
            renderFolder(this.play_queue, this.folder_path),
            html`<span>Loading ...</span>`
        )}`;
    }
}

customElements.define("mss-folder-details", FolderDetails);
