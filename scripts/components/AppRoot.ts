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

import { LitElement, html, TemplateResult } from "lit-element";
import { router } from "../router";
import "./SideBarLink";
import "./TopLevelFolders";
import "./FoldersList";
import "./FolderCover";

const DEFAULT_PAGE = "default";
const FOLDERS_PAGE = "folders";

class AppRoot extends LitElement {
    private current_page = DEFAULT_PAGE;

    constructor() {
        super();

        router
            .on(() => {
                this.current_page = DEFAULT_PAGE;
                this.requestUpdate();
            })
            .on("/folders", () => {
                this.current_page = FOLDERS_PAGE;
                this.requestUpdate();
            })
            .resolve();
    }

    render(): TemplateResult {
        switch (this.current_page) {
            case FOLDERS_PAGE:
                return html`<mss-top-level-folders></mss-top-level-folders> `;
            case DEFAULT_PAGE:
            default:
                return html`Home`;
        }
    }
}

customElements.define("app-root", AppRoot);
