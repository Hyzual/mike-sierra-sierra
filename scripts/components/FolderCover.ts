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
    TemplateResult,
    PropertyDeclarations,
} from "lit-element";

// By importing the SVG image like this, webpack can hash its filename and put
// it in the assets folder.
import svg from "../../images/assets/no-cover.svg";

class FolderCover extends LitElement {
    private folder_id!: number;
    private folder_title!: string;

    static get properties(): PropertyDeclarations {
        return {
            folder_id: { type: Number },
            folder_title: { type: String },
        };
    }

    static readonly styles = css`
        :host {
            display: flex;
            flex-direction: column;
        }

        .folder-link {
            flex: 1 0 auto;
            border: 1px solid var(--dark-accent-color);
            background: var(--dark-shades-color);
            color: var(--light-shades-color);
            text-decoration: none;
        }

        .folder-link:hover {
            border-color: var(--light-accent-color);
            background: var(--lighter-dark-shades-color);
            color: var(--lighter-light-shades-color);
        }

        .folder-header {
            padding: 8px;
        }
    `;

    render(): TemplateResult {
        const folder_uri = `/app/folders/${encodeURIComponent(this.folder_id)}`;

        return html` <a
            href="${folder_uri}"
            class="folder-link"
            title="Browse folder"
        >
            <img src="${svg}" alt="Default cover image" />
            <div class="folder-header">
                <span>${this.folder_title}</span>
            </div>
        </a>`;
    }
}

customElements.define("folder-cover", FolderCover);
