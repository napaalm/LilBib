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

	"git.antonionapolitano.eu/napaalm/LilBib/internal/db"
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
	idParsed, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		// id invalido: torna all'elenco
		http.Redirect(w, r, "/libri/0", http.StatusSeeOther)
		return
	}
	idLibro := uint32(idParsed)
	// TODO manca libro.html!
	libro, err := db.GetLibro(idLibro)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	templates.ExecuteTemplate(w, "libri.html", libro)
}

func HandleLibri(w http.ResponseWriter, r *http.Request) {
	pageStr := strings.TrimPrefix(r.URL.Path, "/libri/")
	page, err := strconv.ParseUint(pageStr, 10, 32)
	if err != nil {
		http.Redirect(w, r, "/libri/0", http.StatusSeeOther)
		return
	}

	q := r.URL.Query()
	titolo := q.Get("titolo")
	nomeAutore := q.Get("autore")
	nomeGenere := q.Get("genere")

	autori, err := db.RicercaAutori(nomeAutore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// senz'offesa, questa interfaccia mi fa un po' schifo
	idsAutori := make([]uint32, 0, len(autori))
	for _, a := range autori {
		idsAutori = append(idsAutori, a.Codice)
	}
	generi, err := db.RicercaGeneri(nomeGenere)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	idsGeneri := make([]uint32, 0, len(generi))
	for _, g := range generi {
		idsGeneri = append(idsGeneri, g.Codice)
	}
	libri, err := db.RicercaLibri(titolo, idsAutori, idsGeneri)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templates.ExecuteTemplate(w, "libri.html", struct {
		Pagina uint64
		Libri  []db.Libro
	}{page, libri})
}
func HandleAutori(w http.ResponseWriter, r *http.Request) {
	initStr := strings.TrimPrefix(r.URL.Path, "/autori/")
	if len(initStr) != 1 {
		http.Redirect(w, r, "/autori/a", http.StatusSeeOther)
		return
	}
	initial := initStr[0]
	if 'A' <= initial && initial <= 'Z' {
		http.Redirect(w, r, "/autori/"+string(initial-'A'+'a'), http.StatusSeeOther)
		return
	}
	if !('a' <= initial && initial <= 'z') {
		http.Redirect(w, r, "/autori/a", http.StatusSeeOther)
		return
	}

	autori, err := db.GetAutori(initial)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	templates.ExecuteTemplate(w, "autori.html", struct {
		Iniziale byte
		Autori   []db.Autore
	}{initial, autori})
}
func HandleGeneri(w http.ResponseWriter, r *http.Request) {
	generi, err := db.GetGeneri()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	templates.ExecuteTemplate(w, "generi.html", struct {
		Generi []db.Genere
	}{generi})
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
