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
    html,
    LitElement,
    PropertyDeclarations,
    TemplateResult,
    css,
} from "lit-element";
import { router } from "../router";

export class SideBarLink extends LitElement {
    private uri!: string;
    private label!: string;

    static get properties(): PropertyDeclarations {
        return {
            uri: { type: String },
            label: { type: String },
        };
    }

    static readonly styles = css`
        .link {
            display: block;
            padding: 4px 8px;
            border-left: 4px solid transparent;
            color: var(--dark-accent-color);
        }

        .link:active,
        .link:hover {
            border-left-color: var(--light-accent-color);
            color: var(--lighter-dark-accent-color);
        }
    `;

    render(): TemplateResult {
        return html`<a href="${this.uri}" @click="${this.navigate}" class="link"
            ><slot name="icon"></slot></i>${this.label}</a
        >`;
    }

    private navigate(event: Event): void {
        event.preventDefault();
        router.navigate(this.uri);
    }
}

customElements.define("side-bar-link", SideBarLink);