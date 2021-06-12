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

import type { Hybrids } from "hybrids";
import { define, html } from "hybrids";
import { router } from "../../router";

const getFullURI = (uri: string): string => router.link(uri);

export interface SidebarLinkElement {
    uri: string;
    label: string;
}

export const navigate = (host: SidebarLinkElement, event: Event): void => {
    event.preventDefault();
    router.navigate(host.uri);
};

export const SidebarLink: Hybrids<SidebarLinkElement> = {
    uri: "",
    label: "",
    render: (host) => {
        return html`<style>
                a {
                    display: block;
                    padding: 4px 8px;
                    border-left: 4px solid transparent;
                    color: var(--dark-accent-color);
                }

                a:focus,
                a:active,
                a:hover {
                    border-left-color: var(--light-accent-color);
                    color: var(--lighter-dark-accent-color);
                }</style
            ><a href="${getFullURI(host.uri)}" onclick="${navigate}"
                ><slot></slot>${host.label}</a
            >`;
    },
};

define("mss-sidebar-link", SidebarLink);
