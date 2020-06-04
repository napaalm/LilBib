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

	"git.antonionapolitano.eu/napaalm/LilBib/internal/auth"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/config"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/db"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/handlers"
)

var Version string

func main() {
	if err := config.LoadConfig("config/config.toml"); err != nil {
		fmt.Println(err)
		return
	}

	if err := db.InizializzaDB(); err != nil {
		fmt.Println(err)
		return
	}

	auth.InitializeSigning()

	fmt.Println("LilBib versione: " + Version)

	// Imposta la versione nei package che lo richiedono
	handlers.Version = Version

	mux := http.NewServeMux()

	// I pattern che finiscono per '/' comprendono anche i sottopercorsi.
	// Sono valutati 'a partire dal più specifico', quindi '/' sarà
	// sempre l'ultimo.
	mux.HandleFunc("/", handlers.HandleRootOr404)
	mux.HandleFunc("/libri/", handlers.HandleLibri)
	mux.HandleFunc("/libro/", handlers.HandleLibro)
	mux.HandleFunc("/autori/", handlers.HandleAutori)
	mux.HandleFunc("/generi", handlers.HandleGeneri)
	mux.HandleFunc("/login", handlers.HandleLogin)
	mux.HandleFunc("/prestito", handlers.HandlePrestito)
	mux.HandleFunc("/restituzione", handlers.HandleRestituzione)
	mux.HandleFunc("/utente", handlers.HandleUtente)
	mux.HandleFunc("/api/getLibro", handlers.HandleGetLibro)
	mux.HandleFunc("/api/prestito", handlers.HandleNewPrestito)
	mux.HandleFunc("/api/restituzione", handlers.HandleSetRestituzione)

	// File server per servire direttamente i contenuti statici.
	fileserver := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileserver))

	srvAddress := config.Config.Generale.Porta
	srv := &http.Server{
		Addr:    srvAddress,
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())
}
