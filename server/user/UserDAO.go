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

package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/swithek/sessionup"
)

// Store handles database operations related to Users
type Store interface {
	GetUserMatchingEmail(ctx context.Context, email string) (*PossibleMatch, error)
	GetUserMatchingSession(ctx context.Context) (*CurrentUser, error)
	SaveFirstAdministrator(ctx context.Context, registration *Registration) error
}

// DAO implements Store
type DAO struct {
	db *sql.DB
}

// NewDAO creates a new DAO
func NewDAO(db *sql.DB) Store {
	return &DAO{db}
}

// GetUserMatchingEmail retrieves the user matching the provided email address.
// If the email is found, it will return a PossibleMatch. If it isn't found, it will return an error.
func (d *DAO) GetUserMatchingEmail(ctx context.Context, email string) (*PossibleMatch, error) {
	query := `SELECT user.id, user.email, user.password FROM user WHERE user.email = ?`
	row := d.db.QueryRowContext(ctx, query, email)
	var (
		id           uint
		passwordHash []byte
	)
	if err := row.Scan(&id, &email, &passwordHash); err != nil {
		return nil, fmt.Errorf("Could not retrieve the user by its credentials: %w", err)
	}

	return &PossibleMatch{id, email, passwordHash}, nil
}

// PossibleMatch represents a user with the same email credential as the one provided
// in the sign-in form. It still needs to check the password matches.
type PossibleMatch struct {
	ID           uint
	Email        string
	PasswordHash []byte
}

// GetUserMatchingSession retrieves the current user from the current session attached to
// the request context.
func (d *DAO) GetUserMatchingSession(ctx context.Context) (*CurrentUser, error) {
	session, _ := sessionup.FromContext(ctx)
	query := `SELECT user.id, user.email, user.username FROM user WHERE user.id = ?`
	row := d.db.QueryRowContext(ctx, query, session.UserKey)
	var (
		id       uint
		email    string
		username string
	)
	if err := row.Scan(&id, &email, &username); err != nil {
		return nil, fmt.Errorf("Could not retrieve current user from session: %w", err)
	}
	return &CurrentUser{id, email, username}, nil
}

// CurrentUser represents the currently signed-in user authentified by its session.
type CurrentUser struct {
	ID       uint
	Email    string
	Username string
}

// SaveFirstAdministrator creates the first administrator account whose
// credentials are given in the registration.
func (d *DAO) SaveFirstAdministrator(ctx context.Context, registration *Registration) error {
	query := `INSERT INTO user(email, password, username) VALUES (?, ?, ?)`
	_, err := d.db.ExecContext(ctx, query, registration.Email, registration.PasswordHash, registration.Username)
	return err
}

// Registration represents the data needed to save the user in databse. Instead of a password,
// it needs a password hash.
type Registration struct {
	Email        string
	PasswordHash []byte
	Username     string
}
