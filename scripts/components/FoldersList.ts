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
import "./FolderCover";

interface Folder {
    id: number;
    title: string;
}

//TODO: This list of folders must come from the REST API
const folders = [
    { id: 1, title: "Ghost in the Shell - Stand Alone Complex OST 3" },
    { id: 2, title: "Call To Power 2" },
    { id: 3, title: "Civilization: Call To Power" },
    { id: 4, title: "Medieval II Total War" },
    { id: 5, title: "Age of Empires Definitive Edition (Original Soundtrack)" },
    { id: 6, title: "Stellaris Digital Soundtrack" },
    { id: 7, title: "WarCraft III: Reign of Chaos [Ripped]" },
    { id: 8, title: "Zeus: Master of Olympus" },
    { id: 9, title: "Il était une fois... l'Homme" },
    {
        id: 10,
        title: "Ghost in the Shell - Stand Alone Complex : Solid State Society",
    },
    { id: 11, title: "Final Fantasy X OST" },
    { id: 12, title: "Video Games Music" },
    { id: 13, title: "Trance" },
    { id: 14, title: "Starcraft OST" },
];

class FoldersList extends LitElement {
    static get properties(): PropertyDeclarations {
        return {
            folders: { type: Array },
        };
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
        return html`${folders.map(
            (folder: Folder) =>
                html`<folder-cover
                    folder_id="${folder.id}"
                    folder_title="${folder.title}"
                ></folder-cover>`
        )}`;
    }
}

customElements.define("folders-list", FoldersList);
