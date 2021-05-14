/*
 *   Copyright (C) 2020-2021  Joris MASSON
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

package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/hyzual/mike-sierra-sierra"
	"github.com/hyzual/mike-sierra-sierra/server/adapter"
	"github.com/hyzual/mike-sierra-sierra/server/adapter/rest"
	"github.com/hyzual/mike-sierra-sierra/server/adapter/server"
	"github.com/hyzual/mike-sierra-sierra/server/adapter/server/app"
	"github.com/hyzual/mike-sierra-sierra/server/adapter/server/user"
	"github.com/hyzual/mike-sierra-sierra/server/domain/music"
	sqlitestore "github.com/hyzual/sessionup-sqlitestore"
	_ "github.com/mattn/go-sqlite3"
	"github.com/swithek/sessionup"
)

const disableHTTPSEnv = "MIKE_DISABLE_HTTPS"

func main() {
	cwd, err := mike.Cwd()
	if err != nil {
		log.Fatalf("could not read the current working directory: %v", err)
	}
	db, err := sql.Open("sqlite3", "file:database/file/mike.db?mode=rwc&cache=shared")
	if err != nil {
		log.Fatalf("could not connect to the database: %v", err)
	}
	db.SetConnMaxLifetime(1 * time.Hour)

	userStore := user.NewDAO(db)
	assetsPath := path.Join(cwd, "assets")
	assetsLoader := adapter.NewBasePathJoiner(cwd)
	templatesPath := path.Join(cwd, "templates")
	templateExecutor := adapter.NewTemplateExecutor(templatesPath)
	musicLoader := adapter.NewBasePathJoiner(music.MusicPath)
	assetsResolver := adapter.NewAssetsResolver(os.DirFS(assetsPath), "/assets")

	sessionStore, err := sqlitestore.New(db, "sessions", time.Minute*30)
	if err != nil {
		log.Fatalf("error while creating a new sessions Store: %v", err)
	}
	sessionManager := sessionup.NewManager(
		sessionStore,
		sessionup.CookieName("id"),
		sessionup.Reject(server.HandleUnauthorized),
	)
	router := mux.NewRouter()
	decoder := schema.NewDecoder()
	user.Register(
		router,
		templateExecutor,
		assetsResolver,
		userStore,
		sessionManager,
		decoder,
	)
	musicDirFS := os.DirFS(music.MusicPath)
	musicLibraryFileSystem := adapter.NewOSFileSystem(musicDirFS, musicLoader)
	explorer := music.NewMusicLibraryExplorer(musicLibraryFileSystem)
	rest.Register(router, sessionManager, explorer)
	app.Register(
		router,
		templateExecutor,
		assetsResolver,
		userStore,
		sessionManager,
	)
	server.Register(
		router,
		sessionManager,
		assetsLoader,
		musicLoader,
	)

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
