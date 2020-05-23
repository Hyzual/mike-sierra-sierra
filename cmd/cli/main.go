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

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hyzual/mike-sierra-sierra/scanner"
)

const musicPath = "/music" // It is a volume in the Docker image

func main() {
	version := flag.Bool("version", false, "Show the program version")

	flag.Parse()

	if *version == true {
		fmt.Println("Mike-Sierra-Sierra CLI v0.1.0")
		return
	}

	arg := flag.Arg(0)
	err := scanner.ScanMusicFiles(arg)
	if err != nil {
		log.Fatalf("Error while scanning music files: %v", err)
	}
}
