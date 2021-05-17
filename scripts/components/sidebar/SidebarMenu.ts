/*
 *   Copyright (C) 2021  Joris MASSON
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

import type { TemplateResult } from "lit";
import { css, html, LitElement, svg } from "lit";
import { unsafeSVG } from "lit/directives/unsafe-svg";
import home_icon from "@glyphs/core/duo/house.svg";
import music_note_icon from "@glyphs/core/duo/music-note.svg?raw";

const scaleIcon = (icon_svg: string): string =>
    icon_svg
        .replace(/width="80"/, 'width="24"')
        .replace(/height="80"/, 'height="24"');

class SidebarMenu extends LitElement {
    static readonly styles = css`
        .title {
            margin: 8px 0 0 12px;
            font-size: 20px;
        }

        .menu {
            margin: 8px 0 0;
            padding: 0;
            list-style: none;
        }

        .icon {
            position: relative;
            top: 6px;
            width: 24px;
            margin: 0 4px 0 0;
        }
    `;

    render(): TemplateResult {
        // const home_icon_svg = svg`${unsafeSVG(home_icon)}`;
        const music_note_svg = svg`${unsafeSVG(scaleIcon(music_note_icon))}`;

        return html`<section>
            <h3 class="title">Music</h3>
            <ul class="menu">
                <li>
                    <mss-side-bar-link uri="/" label="Home">
                        <img
                            class="icon"
                            src="${home_icon}"
                            aria-hidden="true"
                        />
                    </mss-side-bar-link>
                </li>
                <li>
                    <mss-side-bar-link uri="/" label="All Songs">
                        ${music_note_svg}
                    </mss-side-bar-link>
                </li>
                <li>
                    <mss-side-bar-link uri="folders" label="Browse by Folder">
                        <i
                            class="fa fa-fw fa-folder mss-button-icon"
                            aria-hidden="true"
                            slot="icon"
                        ></i>
                    </mss-side-bar-link>
                </li>
                <li>
                    <mss-side-bar-link uri="/" label="Albums">
                        <i
                            class="fa fa-fw fa-picture-o mss-button-icon"
                            aria-hidden="true"
                            slot="icon"
                        ></i>
                    </mss-side-bar-link>
                </li>
                <li>
                    <mss-side-bar-link uri="/" label="Artists">
                        <i
                            class="fa fa-fw fa-user-circle mss-button-icon"
                            aria-hidden="true"
                            slot="icon"
                        ></i>
                    </mss-side-bar-link>
                </li>
                <li>
                    <mss-side-bar-link uri="/" label="Genres">
                        <i
                            class="fa fa-fw fa-tags mss-button-icon"
                            aria-hidden="true"
                            slot="icon"
                        ></i>
                    </mss-side-bar-link>
                </li>
            </ul>
        </section> `;
    }
}

customElements.define("mss-sidebar-menu", SidebarMenu);
