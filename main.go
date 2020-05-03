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
	"log"
	"net/http"

	"github.com/hyzual/mike-sierra-sierra/server"
)

func main() {
	cwd, err := cwd()
	if err != nil {
		log.Fatalf("could not read the current working directory %v", err)
	}

	server := server.New(cwd)
	err = http.ListenAndServe(":8080", server)
	if err != nil {
		log.Fatalf("could not listen on port 8080 %v", err)
	}
}
