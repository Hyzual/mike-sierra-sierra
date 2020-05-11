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

package server

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/pkg/errors"
)

// Borrowed from https://stackoverflow.com/a/32351471

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random number generator fails
// to function correctly, in which case the caller should not continue.
func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, errors.Wrap(err, "Could not generate random bytes")
	}
	return bytes, nil
}

// GenerateRandomString returns a URL-safe base64 encoded securely generated random string.
// length is the number of random bytes. It will NOT match the length of the string as it is base64 encoded.
// It will return an error if the system's secure random number generator fails
// to function correctly, in which case the caller should not continue.
func GenerateRandomString(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	return base64.URLEncoding.EncodeToString(bytes), errors.Wrap(err, "Could not generate random string")
}
