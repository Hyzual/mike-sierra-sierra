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

interface Icon {
    src: string;
}

const Icon: Hybrids<Icon> = {
    src: "",
    render: (host) => {
        return html`<style>
                svg {
                    width: 24px;
                    height: 24px;
                }
                path {
                    stroke: var(--dark-accent-color);
                    stroke-width: 3px;
                }</style
            ><i innerHTML="${host.src}"></i>`;
    },
};

define("mss-icon", Icon);
