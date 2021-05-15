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

import type { PropertyDeclarations, TemplateResult } from "lit-element";
import { LitElement, css, html } from "lit-element";
import type { SubFolder } from "../types";

class FoldersList extends LitElement {
    folders: SubFolder[] = [];

    static get properties(): PropertyDeclarations {
        return { folders: { type: Array } };
    }

    static readonly styles = css`
        :host {
            display: grid;
            gap: 8px 8px;
            grid-auto-flow: row;
            grid-template-columns: repeat(auto-fit, 256px);
            grid-template-rows: max-content;
        }
    `;

    render(): TemplateResult {
        return html`${this.folders.map(
            (folder: SubFolder) =>
                html`<mss-folder-cover
                    folder_title="${folder.name}"
                    folder_uri="${folder.uri}"
                ></mss-folder-cover>`
        )}`;
    }
}

customElements.define("mss-folders-list", FoldersList);
