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

import { css, html, LitElement, TemplateResult } from "lit-element";
import { until } from "lit-html/directives/until";
import { NetworkError } from "../api/NetworkError";
import { getTopFolders } from "../api/rest-querier";

const TopLevelLoadingErrorState = (
    error: Error | NetworkError
): TemplateResult => {
    const error_template =
        error instanceof NetworkError
            ? html`<p class="error">
                  An error occurred: Code ${error.statusCode}:
                  ${error.statusText}
              </p>`
            : html`<p class="error">An error occurred: ${error.message}</p>`;
    return html`<div class="error-container">${error_template}</div>`;
};

const renderTopLevelFolders = async (): Promise<TemplateResult> => {
    const result = await getTopFolders();
    if (result.isErr()) {
        return TopLevelLoadingErrorState(result.error);
    }
    return html`<folders-list .folders=${result.value.folders}></folders-list>`;
};

class TopLevelFoldersElement extends LitElement {
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
            renderTopLevelFolders(),
            html`<span>Loading ...</span>`
        )}`;
    }
}

customElements.define("mss-top-level-folders", TopLevelFoldersElement);
