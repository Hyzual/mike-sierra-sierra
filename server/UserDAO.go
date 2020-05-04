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
	"database/sql"

	"github.com/pkg/errors"
)

// UserStore handles database operations linked to Users
type UserStore interface {
	VerifyCredentialsMatch(credentials LoginFormRepresentation) error
}

type UserDAO struct {
	db *sql.DB
}

func NewUserDao(db *sql.DB) UserStore {
	return &UserDAO{db}
}

func (u *UserDAO) VerifyCredentialsMatch(credentials LoginFormRepresentation) error {
	row := u.db.QueryRow(`SELECT * FROM user
		WHERE user.email = ? AND user.password = ?`, credentials.Email, credentials.Password)

	var (
		id       uint64
		email    string
		password string
	)
	if err := row.Scan(&id, &email, &password); err != nil {
		return errors.Wrap(err, "Could not retrieve the user by its credentials")
	}

	return nil
}
