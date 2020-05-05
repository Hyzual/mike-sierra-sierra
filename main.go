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

package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/hyzual/mike-sierra-sierra/server"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cwd, err := cwd()
	if err != nil {
		log.Fatalf("could not read the current working directory %v", err)
	}
	db, err := sql.Open("sqlite3", "file:database/file/mike.db?mode=rwc&cache=shared")
	if err != nil {
		log.Fatalf("could not connect to the database %v", err)
	}
	db.SetConnMaxLifetime(1 * time.Hour)
	// err = db.Ping()
	// if err != nil {
	// 	log.Fatalf("could not ping the database %v", err)
	// }

	userStore := server.NewUserDao(db)
	pathJoiner := server.NewRootPathJoiner(cwd)
	loginHandler := server.NewLoginHandler(pathJoiner, userStore)
	router := server.New(pathJoiner, loginHandler)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("could not listen on port 8080 %v", err)
	}
}
