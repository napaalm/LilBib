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
	"fmt"
	"log"
	"net/http"
)

// TODO config fatto bene
const srvAddress = ":8081"

func rootOr404(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintln(w, "It works!")
}
func libri(w http.ResponseWriter, r *http.Request)        {}
func libro(w http.ResponseWriter, r *http.Request)        {}
func autori(w http.ResponseWriter, r *http.Request)       {}
func generi(w http.ResponseWriter, r *http.Request)       {}
func login(w http.ResponseWriter, r *http.Request)        {}
func prestito(w http.ResponseWriter, r *http.Request)     {}
func restituzione(w http.ResponseWriter, r *http.Request) {}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", rootOr404)
	mux.HandleFunc("/libri/", libri)
	mux.HandleFunc("/libro/", libro)
	mux.HandleFunc("/autori/", autori)
	mux.HandleFunc("/generi/", generi)
	mux.HandleFunc("/login/", login)
	mux.HandleFunc("/prestito/", prestito)
	mux.HandleFunc("/restituzione/", restituzione)

	fileserver := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileserver))

	srv := &http.Server{
		Addr:    srvAddress,
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())

}
