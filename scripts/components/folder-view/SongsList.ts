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

import type { PropertyDeclarations, TemplateResult } from "lit";
import { LitElement, html, css } from "lit";
import type { Song } from "../../types";
import type { PlayQueueState } from "../music/PlayQueueState";

class SongsList extends LitElement {
    songs: Song[] = [];
    play_queue!: PlayQueueState;

    static get properties(): PropertyDeclarations {
        return { songs: { type: Array } };
    }

    static readonly styles = css`
        :host {
            display: flex;
            flex-direction: column;
        }
    `;

    render(): TemplateResult {
        return html`${this.songs.map(
            (song: Song) =>
                html`<mss-song-line
                    .song=${song}
                    .play_queue=${this.play_queue}
                ></mss-song-line>`
        )}`;
    }
}

customElements.define("mss-songs-list", SongsList);
