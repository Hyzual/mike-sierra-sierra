/*
 *   Copyright (c) 2020 Joris MASSON

 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.

 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.

 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package server_test

import (
	"testing"

	"github.com/hyzual/mike-sierra-sierra/server"
)

func TestRandomBytes(t *testing.T) {
	rand, err := server.GenerateRandomBytes(32)

	assertNoError(t, err)
	assertLengthEquals(t, len(rand), 32)
}

func TestRandomString(t *testing.T) {
	rand, err := server.GenerateRandomString(32)

	assertNoError(t, err)
	assertLengthEquals(t, len(rand), 44)
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("did not expect an error but got one %v", err)
	}
}

func assertLengthEquals(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("got length of %d, want %d", got, want)
	}
}
