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

import type { Hybrids } from "hybrids";
import { define, html } from "hybrids";
import home_icon from "@glyphs/core/path/house.svg?raw";
import music_note_icon from "@glyphs/core/path/music-note.svg?raw";
import folder_open_icon from "@glyphs/core/path/folder-open.svg?raw";
import image_one_icon from "@glyphs/core/path/image-1.svg?raw";
import user_circle_icon from "@glyphs/core/path/user-circle.svg?raw";
import tags_icon from "@glyphs/core/path/tags.svg?raw";

type SidebarMenu = Record<never, string>;

const SidebarMenu: Hybrids<SidebarMenu> = {
    render: () => {
        return html`<section>
                <h3 class="title">Music</h3>
                <ul class="menu">
                    <li>
                        <mss-sidebar-link uri="/" label="Home">
                            <mss-icon
                                class="icon"
                                src="${home_icon}"
                            ></mss-icon>
                        </mss-sidebar-link>
                    </li>
                    <li>
                        <mss-sidebar-link uri="/" label="All Songs">
                            <mss-icon
                                class="icon"
                                src="${music_note_icon}"
                            ></mss-icon>
                        </mss-sidebar-link>
                    </li>
                    <li>
                        <mss-sidebar-link
                            uri="folders"
                            label="Browse by Folder"
                        >
                            <mss-icon
                                class="icon"
                                src="${folder_open_icon}"
                            ></mss-icon>
                        </mss-sidebar-link>
                    </li>
                    <li>
                        <mss-sidebar-link uri="/" label="Albums">
                            <mss-icon
                                class="icon"
                                src="${image_one_icon}"
                            ></mss-icon>
                        </mss-sidebar-link>
                    </li>
                    <li>
                        <mss-sidebar-link uri="/" label="Artists">
                            <mss-icon
                                class="icon"
                                src="${user_circle_icon}"
                            ></mss-icon>
                        </mss-sidebar-link>
                    </li>
                    <li>
                        <mss-sidebar-link uri="/" label="Genres">
                            <mss-icon
                                class="icon"
                                src="${tags_icon}"
                            ></mss-icon>
                        </mss-sidebar-link>
                    </li>
                </ul>
            </section>
            <style>
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
                    margin: 0 4px 0 0;
                }
            </style>`;
    },
};

define("mss-sidebar-menu", SidebarMenu);
