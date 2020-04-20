/*
 * handlers.go
 *
 * Descrizione del file.
 *
 * Copyright (c) 2020 Nome Cognome <nome.cognome@example.org>
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

package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

const templatesDir = "web/template"

var templates = template.Must(template.ParseFiles(
	templatesDir+"/autori.html",
	templatesDir+"/generi.html",
	templatesDir+"/index.html",
	templatesDir+"/libri.html",
	templatesDir+"/login.html",
	templatesDir+"/prestito.html",
	templatesDir+"/restituzione.html",
))

func HandleRootOr404(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	HandleHome(w, r)
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func HandleLibro(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/libro/")
	_, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		// id invalido: torna all'elenco
		http.Redirect(w, r, "/libri/0", 303)
	}
	// TODO manca libro.html!
	templates.ExecuteTemplate(w, "libri.html", nil)
}

func HandleLibri(w http.ResponseWriter, r *http.Request) {
	pageStr := strings.TrimPrefix(r.URL.Path, "/libri/")
	_, err := strconv.ParseUint(pageStr, 10, 32)
	if err != nil {
		// pagina non valida
		http.Redirect(w, r, "/libri/0", 303)
	}

	// q := r.URL.Query()
	// titolo := q.Get("titolo")
	// autore := q.Get("autore")
	// genere := q.Get("genere")

	// idsAutore := db.RicercaAutori(autore)
	// idsGeneri := db.RicercaGeneri(genere)
	templates.ExecuteTemplate(w, "libri.html", nil)
}
func HandleAutori(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/autori/" {
		// HTTP 303 See Other
		http.Redirect(w, r, "/autori/a", 303)
		return
	}
	templates.ExecuteTemplate(w, "autori.html", nil)
}
func HandleGeneri(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "generi.html", nil)
}
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}
func HandleUtente(w http.ResponseWriter, r *http.Request) {
	// TODO manca utente!
	templates.ExecuteTemplate(w, "login.html", nil)
}
func HandlePrestito(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "prestito.html", nil)
}
func HandleRestituzione(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "restituzione.html", nil)
}
