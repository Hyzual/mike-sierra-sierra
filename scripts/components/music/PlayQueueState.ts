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

import type { Song } from "../../types";
import { NullSong } from "../../types";

type CurrentSongChangedCallback = (new_song: Song) => void;

const noopSongChangedCallback = (): void => {
    //Do nothing
};

export class PlayQueueState {
    private callback: CurrentSongChangedCallback = noopSongChangedCallback;
    #currentSong: Song = NullSong;

    get currentSong(): Song {
        return this.#currentSong;
    }

    set currentSong(new_song: Song) {
        this.callback(new_song);
        this.#currentSong = new_song;
    }

    setCallback(callback: CurrentSongChangedCallback): void {
        this.callback = callback;
    }
}
