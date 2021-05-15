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

import type { PropertyDeclarations, TemplateResult } from "lit-element";
import { LitElement, html, css } from "lit-element";
import type { Song } from "scripts/types";
import type { PlayQueueState } from "../music/PlayQueueState";

export class SongLine extends LitElement {
    song!: Song;
    play_queue!: PlayQueueState;

    static get properties(): PropertyDeclarations {
        return { song: { type: { Object } } };
    }

    static readonly styles = css`
        :host {
            display: flex;
            flex: 1;
        }
    `;

    render(): TemplateResult {
        return html`<span>${this.song.title}</span>
            <a href="${this.song.uri}" @click="${this.playSong}">Play</a>`;
    }

    private playSong(event: Event): void {
        event.preventDefault();
        this.play_queue.currentSong = this.song;
    }
}

customElements.define("mss-song-line", SongLine);
