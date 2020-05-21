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

package user

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

// Store handles database operations related to Users
type Store interface {
	GetUserMatchingCredentials(ctx context.Context, credentials *Credentials) (*Current, error)
	SaveFirstAdministrator(ctx context.Context, form *RegistrationForm) error
}

// DAO implements Store
type DAO struct {
	db *sql.DB
}

// NewDAO creates a new DAO
func NewDAO(db *sql.DB) Store {
	return &DAO{db}
}

// GetUserMatchingCredentials retrieves the user matching the provided credentials.
// If the credentials match credentials stored in the database, it will return a Current and nil error.
// If they don't match, it will return an error
func (d *DAO) GetUserMatchingCredentials(ctx context.Context, credentials *Credentials) (*Current, error) {
	query := `SELECT user.id, user.email FROM user WHERE user.email = ? AND user.password = ?`
	row := d.db.QueryRowContext(ctx, query, credentials.Email, credentials.Password)
	var (
		id    uint
		email string
	)
	if err := row.Scan(&id, &email); err != nil {
		return nil, errors.Wrap(err, "Could not retrieve the current user by its credentials")
	}

	return &Current{id, email}, nil
}

// SaveFirstAdministrator creates the first administrator account whose
// credentials are given in the registration form.
func (d *DAO) SaveFirstAdministrator(ctx context.Context, form *RegistrationForm) error {
	query := `INSERT INTO user(email, password, username) VALUES (?, ?, ?)`
	_, err := d.db.ExecContext(ctx, query, form.Email, form.Password, form.Username)
	return err
}
