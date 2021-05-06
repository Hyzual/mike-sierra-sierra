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

import {
    css,
    html,
    LitElement,
    PropertyDeclarations,
    TemplateResult,
} from "lit-element";
import { until } from "lit-html/directives/until";
import { NetworkError } from "../api/NetworkError";
import { getFolder } from "../api/rest-querier";

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

const renderFolder = async (folder_path: string): Promise<TemplateResult> => {
    const result = await getFolder(folder_path);
    if (result.isErr()) {
        return renderErrorState(result.error);
    }
    return html`<folders-list .folders=${result.value.folders}></folders-list>`;
};

export class FolderDetails extends LitElement {
    folder_path = "";

    static get properties(): PropertyDeclarations {
        return { folder_path: { type: Object } };
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
            renderFolder(this.folder_path),
            html`<span>Loading ...</span>`
        )}`;
    }
}

customElements.define("mss-folder-details", FolderDetails);
