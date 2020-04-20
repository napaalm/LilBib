/*
 * main.go
 *
 * Codice principale del programma.
 *
 * Copyright (c) 2020 Antonio Napolitano <nap@antonionapolitano.eu>
 *
 * This file is part of LilBib.
 *
 * LilBib is free software; you can redistribute it and/or modify it
 * under the terms of the Affero GNU General Public License as
 * published by the Free Software Foundation; either version 3, or (at
 * your option) any later version.
 *
 * LilBib is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
 * or FITNESS FOR A PARTICULAR PURPOSE.  See the Affero GNU General
 * Public License for more details.
 *
 * You should have received a copy of the Affero GNU General Public
 * License along with LilBib; see the file LICENSE. If not see
 * <http://www.gnu.org/licenses/>.
 */

package main

import (
	"log"
	"net/http"

	"git.antonionapolitano.eu/napaalm/LilBib/internal/handlers"
)

// TODO config fatto bene
const srvAddress = ":8081"

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HandleRootOr404)
	mux.HandleFunc("/libri/", handlers.HandleLibri)
	mux.HandleFunc("/libro/", handlers.HandleLibri)
	mux.HandleFunc("/autori/", handlers.HandleAutori)
	mux.HandleFunc("/generi/", handlers.HandleGeneri)
	mux.HandleFunc("/login/", handlers.HandleLogin)
	mux.HandleFunc("/prestito/", handlers.HandlePrestito)
	mux.HandleFunc("/restituzione/", handlers.HandleRestituzione)

	fileserver := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileserver))

	srv := &http.Server{
		Addr:    srvAddress,
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())

}
