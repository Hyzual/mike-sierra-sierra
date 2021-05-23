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
import { LitElement, html, css } from "lit";
import { router } from "../router";
import "./sidebar/SidebarMenu";
import "./sidebar/SideBarLink";
import "./sidebar/FolderIcon";
import "./folder-view/FolderDetails";
import "./folder-view/FoldersList";
import "./folder-view/FolderCover";
import "./folder-view/SongsList";
import "./folder-view/SongLine";
import "./music/MusicPlayer";
import { PlayQueueState } from "./music/PlayQueueState";

type Page = "default" | "folders" | "play_queue";
const DEFAULT_PAGE: Page = "default";
const FOLDERS_PAGE: Page = "folders";
const PLAY_QUEUE_PAGE: Page = "play_queue";

class AppRoot extends LitElement {
    private current_page: Page = DEFAULT_PAGE;
    private current_folder_path = "";
    private play_queue: PlayQueueState;

    constructor() {
        super();
        this.play_queue = new PlayQueueState();

        router
            .on(() => {
                this.current_page = DEFAULT_PAGE;
            })
            .on("/play-queue", () => {
                this.current_page = PLAY_QUEUE_PAGE;
            })
            .on("/folders", () => {
                this.current_page = FOLDERS_PAGE;
                this.current_folder_path = "";
            })
            .on("/folders/:path", (match) => {
                this.current_page = FOLDERS_PAGE;
                if (match && match.data) {
                    this.current_folder_path = match.data.path;
                }
            })
            .resolve();
    }

    static get properties(): PropertyDeclarations {
        return {
            current_page: { state: true },
            current_folder_path: { state: true },
        };
    }

    static readonly styles = css`
        :host {
            display: grid;
            grid-template-areas:
                "navbar navbar"
                "sidebar breadcrumbs"
                "sidebar main"
                "footer footer";
            grid-template-columns: minmax(300px, auto) 1fr;
            grid-template-rows: 40px 40px 1fr 40px;
            height: 100vh;
        }

        .main {
            grid-area: main;
            background: var(--dark-shades-color);
        }

        .breadcrumbs {
            grid-area: breadcrumbs;
            background: var(--darker-dark-shades-color);
        }

        .sidebar {
            grid-area: sidebar;
            background: var(--darker-dark-shades-color);
        }

        .footer {
            grid-area: footer;
            background: var(--darker-dark-shades-color);
        }
    `;

    render(): TemplateResult {
        return html`<slot name="header"></slot>
            <nav class="breadcrumbs">Breadcrumbs</nav>
            <mss-sidebar-menu class="sidebar"></mss-sidebar-menu>
            <main class="main">${this.renderMainElement()}</main>
            <mss-music-player
                class="footer"
                .play_queue=${this.play_queue}
            ></mss-music-player>`;
    }

    private renderMainElement(): TemplateResult {
        switch (this.current_page) {
            case FOLDERS_PAGE:
                return html`<mss-folder-details
                    .folder_path=${this.current_folder_path}
                    .play_queue=${this.play_queue}
                ></mss-folder-details> `;
            case DEFAULT_PAGE:
            default:
                return html`Home`;
        }
    }
}

customElements.define("mss-app-root", AppRoot);
