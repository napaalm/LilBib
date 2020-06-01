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

type homeVars struct {
	Disponibili int
	Prenotati   int
}

// viene inizializzato nel momento in cui viene importato il package
var templates = template.Must(template.ParseFiles(
	templatesDir+"/autori.html",
	templatesDir+"/generi.html",
	templatesDir+"/index.html",
	templatesDir+"/libri.html",
	templatesDir+"/login.html",
	templatesDir+"/prestito.html",
	templatesDir+"/restituzione.html",
))

// Handler per qualunque percorso diverso da tutti gli altri percorsi riconosciuti.
// Caso particolare è la homepage (/); per ogni altro restituisce 404.
func HandleRootOr404(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	HandleHome(w, r)
}

// Percorso: /
// Homepage.
func HandleHome(w http.ResponseWriter, r *http.Request) {
	pren, err := db.LibriPrenotati()
	if err != nil {
		//errore, imposto dei valori di default
		templates.ExecuteTemplate(w, "index.html", homeVars{Disponibili: -1, Prenotati: -1})
		return
	}

	disp, err := db.LibriDisponibili()
	if err != nil {
		//errore, imposto dei valori di default
		templates.ExecuteTemplate(w, "index.html", homeVars{Disponibili: -1, Prenotati: -1})
		return
	}

	vars := homeVars{Disponibili: disp, Prenotati: pren}
	templates.ExecuteTemplate(w, "index.html", vars)
}

// Percorso: /libro/<idLibro uint32>
// Mostra informazioni sul libro con codice `idLibro`
// Reindirizza a /libri/0 (elenco dei libri) nel caso di `idLibro` assente o invalido
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

// Formato: /libri/<page uint32>
// Mostra la pagina `page` dei risultati della ricera determinata dalla query GET
// Reindirizza a /libri/0  nel caso di `page` assente o invalida
func HandleLibri(w http.ResponseWriter, r *http.Request) {
	pageStr := strings.TrimPrefix(r.URL.Path, "/libri/")
	pageParsed, err := strconv.ParseUint(pageStr, 10, 32)
	if err != nil {
		http.Redirect(w, r, "/libri/0", http.StatusSeeOther)
		return
	}
	page := uint16(pageParsed)

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
	libri, err := db.RicercaLibri(titolo, idsAutori, idsGeneri, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templates.ExecuteTemplate(w, "libri.html", struct {
		Pagina uint16
		Libri  []db.Libro
	}{page, libri})
}

// Percorso: /autori/<iniziale byte>
// Mostra l'elenco degli autori con iniziale `iniziale`
// Reindirizza a /autori/a nel caso di `iniziale` assente o invalida
func HandleAutori(w http.ResponseWriter, r *http.Request) {
	initStr := strings.TrimPrefix(r.URL.Path, "/autori/")
	if len(initStr) != 1 {
		http.Redirect(w, r, "/autori/a", http.StatusSeeOther)
		return
	}
	initial := initStr[0]
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

// Percorso: /generi
// Mostra l'elenco dei generi.
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

// Percorso: /login
// Mostra pagina di accesso.
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

// Percorso: /utente
// Mostra informazioni sull'utente.
func HandleUtente(w http.ResponseWriter, r *http.Request) {
	// TODO manca utente!
	templates.ExecuteTemplate(w, "login.html", nil)
}

// Percorso: /prestito
// Permette di scansionare o inserire il codice di uno o più libri per prenderli in prestito scegliendone la durata.
func HandlePrestito(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "prestito.html", nil)
}

// Percorso: /restituzione
// Permette di restituire i libri in proprio possesso.
func HandleRestituzione(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "restituzione.html", nil)
}
