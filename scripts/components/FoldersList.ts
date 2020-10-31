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

import {
    LitElement,
    css,
    html,
    PropertyDeclarations,
    TemplateResult,
} from "lit-element";
import { SubFolder } from "../types";

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
                html`<folder-cover
                    folder_id="${folder.id}"
                    folder_title="${folder.name}"
                ></folder-cover>`
        )}`;
    }
}

customElements.define("folders-list", FoldersList);
