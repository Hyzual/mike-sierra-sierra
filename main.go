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
	"os"
	"path"
	"time"

	"github.com/hyzual/mike-sierra-sierra/server"
	"github.com/hyzual/mike-sierra-sierra/server/user"
	_ "github.com/mattn/go-sqlite3"
	"github.com/swithek/sessionup"
	"github.com/swithek/sessionup/memstore"
)

const disableHTTPSEnv = "MIKE_DISABLE_HTTPS"
const musicPath = "/music" // It is a volume in the Docker image

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

	userStore := user.NewDAO(db)
	assetsIncluder := server.NewAssetsIncluder(cwd)
	templatesPath := path.Join(cwd, "templates")
	templateLoader := server.NewTemplateLoader(templatesPath)
	musicLoader := server.NewMusicLoader(musicPath)

	sessionStore := memstore.New(time.Minute * 5)
	sessionManager := sessionup.NewManager(
		sessionStore, sessionup.CookieName("id"),
		sessionup.Reject(server.HandleUnauthorized),
	)
	loginHandler := user.NewLoginHandler(userStore, sessionManager)
	router := server.New(sessionManager, assetsIncluder, templateLoader, musicLoader, loginHandler)

	var port string
	isHTTPSDisabled := os.Getenv(disableHTTPSEnv) != ""
	if isHTTPSDisabled {
		port = "8080"
	} else {
		port = "8443"
	}
	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	if isHTTPSDisabled {
		err = srv.ListenAndServe()
	} else {
		err = srv.ListenAndServeTLS("./secrets/cert.pem", "./secrets/key.pem")
	}
	if err != nil {
		log.Fatalf("could not listen on port %s %v", port, err)
	}
}
