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
import { css, html, LitElement } from "lit";
import type { PlayQueueState } from "./PlayQueueState";

export class MusicPlayer extends LitElement {
    readonly play_queue!: PlayQueueState;

    static get properties(): PropertyDeclarations {
        return { play_queue: { type: Object } };
    }

    firstUpdated(): void {
        this.play_queue.setCallback(this.onCurrentSongChange.bind(this));
    }

    static readonly styles = css`
        :host {
            display: flex;
        }

        .player {
            flex: 1;
        }
    `;

    render(): TemplateResult {
        return html`<audio
            controls
            class="player"
            src="${this.play_queue.currentSong.uri}"
        ></audio>`;
    }

    private onCurrentSongChange(): void {
        this.requestUpdate();
    }
}

customElements.define("mss-music-player", MusicPlayer);
