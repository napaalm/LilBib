/*
 * handlers.go
 *
 * Package per gestire le diverse pagine ed i relativi template.
 *
 * Copyright (c) 2020 Antonio Napolitano <nap@antonionapolitano.eu>
 * Copyright (c) 2020 Davide Vendramin <natalianatiche02@gmail.com>
 * Copyright (c) 2020 Maxim Kovalkov <kov1maxim1al@gmail.com>
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

// Package per gestire le diverse pagine ed i relativi template.
package handlers

import (
	"encoding/json"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/auth"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/config"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/db"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/hash"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/qrcode"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

const templatesDir = "web/template"

var Version string

type CommonValues struct {
	Version string
}

type PrestitoTitolo struct {
	Prestito db.Prestito
	Titolo   string
}

// viene inizializzato nel momento in cui viene importato il package
var templates = template.Must(template.ParseFiles(
	templatesDir+"/autori.html",
	templatesDir+"/generi.html",
	templatesDir+"/index.html",
	templatesDir+"/libri.html",
	templatesDir+"/libro.html",
	templatesDir+"/login.html",
	templatesDir+"/prestito.html",
	templatesDir+"/restituzione.html",
	templatesDir+"/utente.html",
	templatesDir+"/aggiungiLibro.html",
	templatesDir+"/generaCodici.html",
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
		templates.ExecuteTemplate(w, "index.html", struct {
			Disponibili uint32
			Totali      uint32
			Values      CommonValues
		}{0, 0, CommonValues{Version}})
		return
	}

	disp, err := db.LibriDisponibili()
	if err != nil {
		//errore, imposto dei valori di default
		templates.ExecuteTemplate(w, "index.html", struct {
			Disponibili uint32
			Totali      uint32
			Values      CommonValues
		}{0, 0, CommonValues{Version}})
		return
	}

	autoriTot, err := db.CountAutori()
	templates.ExecuteTemplate(w, "index.html", struct {
		Disponibili uint32
		Totali      uint32
		AutoriTot   uint32
		Values      CommonValues
	}{disp, pren + disp, autoriTot, CommonValues{Version}})
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
	libro, err := db.GetLibro(idLibro)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	assegnatario, err := db.GetAssegnatario(idLibro)
	templates.ExecuteTemplate(w, "libro.html", struct {
		Libro  db.Libro
		Utente string
		Values CommonValues
	}{libro, assegnatario, CommonValues{Version}})
}

// Formato: /libri/<page uint32>
// Mostra la pagina `page` dei risultati della ricerca determinata dalla query GET
// Reindirizza a /libri/0  nel caso di `page` assente o invalida
func HandleLibri(w http.ResponseWriter, r *http.Request) {
	pageStr := strings.TrimPrefix(r.URL.Path, "/libri/")
	pageParsed, err := strconv.ParseUint(pageStr, 10, 32)
	if err != nil {
		http.Redirect(w, r, "/libri/0", http.StatusSeeOther)
		return
	}
	page := int16(pageParsed)

	q := r.URL.Query()
	titolo := q.Get("titolo")
	nomeAutore := q.Get("autori")
	nomeGenere := q.Get("generi")
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

	if page == 0 {
		if float64(len(libri))/float64(config.Config.Generale.LunghezzaPagina) <= 1 {
			if float64(page) > (float64(len(libri)) / float64(config.Config.Generale.LunghezzaPagina)) {
				templates.ExecuteTemplate(w, "libri.html", struct {
					PaginaPrec int16
					Pagina     int16
					PaginaSucc int16
					Titolo     string
					Autori     string
					Generi     string
					Libri      []db.Libro
					Values     CommonValues
				}{page, page + 1, page + 1, titolo, nomeAutore, nomeGenere, libri, CommonValues{Version}})
			} else {
				templates.ExecuteTemplate(w, "libri.html", struct {
					PaginaPrec int16
					Pagina     int16
					PaginaSucc int16
					Titolo     string
					Autori     string
					Generi     string
					Libri      []db.Libro
					Values     CommonValues
				}{page, page + 1, page, titolo, nomeAutore, nomeGenere, libri, CommonValues{Version}})
			}
		} else {
			templates.ExecuteTemplate(w, "libri.html", struct {
				PaginaPrec int16
				Pagina     int16
				PaginaSucc int16
				Titolo     string
				Autori     string
				Generi     string
				Libri      []db.Libro
				Values     CommonValues
			}{page, page + 1, page, titolo, nomeAutore, nomeGenere, libri, CommonValues{Version}})
		}

	} else {
		if float64(page) > (float64(len(libri)) / float64(config.Config.Generale.LunghezzaPagina)) {
			templates.ExecuteTemplate(w, "libri.html", struct {
				PaginaPrec int16
				Pagina     int16
				PaginaSucc int16
				Titolo     string
				Autori     string
				Generi     string
				Libri      []db.Libro
				Values     CommonValues
			}{page - 1, page + 1, page, titolo, nomeAutore, nomeGenere, libri, CommonValues{Version}})
		} else {
			templates.ExecuteTemplate(w, "libri.html", struct {
				PaginaPrec int16
				Pagina     int16
				PaginaSucc int16
				Titolo     string
				Autori     string
				Generi     string
				Libri      []db.Libro
				Values     CommonValues
			}{page - 1, page + 1, page + 1, titolo, nomeAutore, nomeGenere, libri, CommonValues{Version}})
		}
	}

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
		Values   CommonValues
	}{initial, autori, CommonValues{Version}})
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
		Values CommonValues
	}{generi, CommonValues{Version}})
}

// Percorso: /login
// Mostra pagina di accesso.
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		username_list, ok0 := r.Form["username"]
		password_list, ok1 := r.Form["password"]
		if !ok0 || !ok1 || len(username_list) != 1 || len(password_list) != 1 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		username := username_list[0]
		password := password_list[0]

		// Controlla le credenziali e ottiene il token
		token, err := auth.AuthenticateUser(username, password)

		// Se l'autenticazione fallisce ritorna 401
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Ottiene il dominio del sito web
		fqdn := config.Config.Generale.FQDN

		// Crea e imposta il cookie
		cookie := http.Cookie{Name: "access_token", Value: string(token), Domain: fqdn, MaxAge: 86400}
		http.SetCookie(w, &cookie)

		// Reindirizza a /utente
		http.Redirect(w, r, "/utente", http.StatusSeeOther)
		return
	}

	// Ottiene il cookie
	_, err := r.Cookie("access_token")

	// Se riesce ad ottenerlo reindirizza ad /utente
	if err == nil {
		http.Redirect(w, r, "/utente", http.StatusSeeOther)
		return
	}

	templates.ExecuteTemplate(w, "login.html", struct {
		Values CommonValues
	}{CommonValues{Version}})
}

// Percorso: /utente
// Mostra informazioni sull'utente.
func HandleUtente(w http.ResponseWriter, r *http.Request) {

	// Ottiene il cookie
	cookie, err := r.Cookie("access_token")

	// Se non riesce ad ottenerlo reindirizza a /login
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Estrae e controlla il token
	token := []byte(cookie.Value)
	utente, err := auth.ParseToken(token)

	// Se l'autenticazione fallisce reindirizza a /login
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	prestiti, err := db.GetPrestiti(utente.Username)

	prestitiTitoli := make([]PrestitoTitolo, len(prestiti))
	for index, prestito := range prestiti {
		libro, err := db.GetLibro(prestito.Libro)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		prestitiTitoli[index].Prestito = prestiti[index]
		prestitiTitoli[index].Titolo = libro.Titolo
	}

	templates.ExecuteTemplate(w, "utente.html", struct {
		Utente         auth.UserInfo
		PrestitiTitoli []PrestitoTitolo
		Values         CommonValues
	}{utente, prestitiTitoli, CommonValues{Version}})
}

// Formato: /api/getLibro?qrcode=<base64-encoded code+password>
// Ritorna le informazioni del libro in formato JSON.
func HandleGetLibro(w http.ResponseWriter, r *http.Request) {

	// Ottiene la password del libro
	q := r.URL.Query()
	password := q.Get("qrcode")

	// Ottiene il libro a partire dalla password
	libro, err := hash.Verifica(password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Ritorna in JSON il libro
	json.NewEncoder(w).Encode(libro)
}

// Percorso: /prestito
// Permette di scansionare o inserire il codice di uno o più libri per prenderli in prestito scegliendone la durata.
func HandlePrestito(w http.ResponseWriter, r *http.Request) {

	// Ottiene il cookie
	cookie, err := r.Cookie("access_token")

	// Se non riesce ad ottenerlo reindirizza a /login
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Estrae e controlla il token
	token := []byte(cookie.Value)
	_, err = auth.ParseToken(token)

	// Se l'autenticazione fallisce reindirizza a /login
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	templates.ExecuteTemplate(w, "prestito.html", struct {
		Values CommonValues
	}{CommonValues{Version}})
}

// Formato: /api/prestito?qrcode=<base64-encoded code+password>&durata=<time in seconds>
// Aggiunge un nuovo prestito per il libro e la durata passati.
func HandleNewPrestito(w http.ResponseWriter, r *http.Request) {

	// Ottiene il cookie
	cookie, err := r.Cookie("access_token")

	// Se non riesce ad ottenerlo ritorna 401
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Estrae e controlla il token
	token := []byte(cookie.Value)
	user, err := auth.ParseToken(token)

	// Se l'autenticazione fallisce ritorna 401
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Ottiene la password del libro e la durata del prestito
	q := r.URL.Query()
	password := q.Get("qrcode")
	durataParsed, err := strconv.ParseUint(q.Get("durata"), 10, 32)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	durata := uint32(durataParsed)

	// Ottiene il libro a partire dalla password
	libro, err := hash.Verifica(password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Ottiene username e id del libro
	username := user.Username
	id := libro.Codice

	// Aggiunge il prestito
	_, err = db.AddPrestito(id, username, durata)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ritorna OK
	http.Error(w, "OK", http.StatusOK)
}

// Percorso: /restituzione
// Permette di restituire i libri in proprio possesso.
func HandleRestituzione(w http.ResponseWriter, r *http.Request) {

	// Ottiene il cookie
	cookie, err := r.Cookie("access_token")

	// Se non riesce ad ottenerlo reindirizza a /login
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Estrae e controlla il token
	token := []byte(cookie.Value)
	_, err = auth.ParseToken(token)

	// Se l'autenticazione fallisce reindirizza a /login
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	templates.ExecuteTemplate(w, "restituzione.html", struct {
		Values CommonValues
	}{CommonValues{Version}})
}

// Formato: /api/restituzione?qrcode=<base64-encoded code+password>
// Imposta come restituito il libro passato in argomento.
func HandleSetRestituzione(w http.ResponseWriter, r *http.Request) {

	// Ottiene il cookie
	cookie, err := r.Cookie("access_token")

	// Se non riesce ad ottenerlo ritorna 401
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Estrae e controlla il token
	token := []byte(cookie.Value)
	_, err = auth.ParseToken(token)

	// Se l'autenticazione fallisce ritorna 401
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Ottiene la password del libro
	q := r.URL.Query()
	password := q.Get("qrcode")

	// Ottiene il libro a partire dalla password
	libro, err := hash.Verifica(password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Ottiene id del libro
	id := libro.Codice

	// Ottiene il prestito corrente
	prestito, err := db.GetCurrentPrestito(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	idPrestito := prestito.Codice

	// Imposta la restituzione
	err = db.SetRestituzione(idPrestito)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ritorna OK
	http.Error(w, "OK", http.StatusOK)
}

func HandleAggiungiLibro(w http.ResponseWriter, r *http.Request) {

	// Ottiene il cookie
	cookie, err := r.Cookie("access_token")

	// Se non riesce ad ottenerlo reindirizza a /login
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Estrae e controlla il token
	token := []byte(cookie.Value)
	user, err := auth.ParseToken(token)

	// Se l'autenticazione fallisce oppure l'utente non è admin ritorna 401
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if !user.IsAdmin {
		http.Error(w, "You are not an admin!", http.StatusUnauthorized)
		return
	}

	generi, err := db.GetGeneri()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	autori, err := db.GetAutori('a')
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	templates.ExecuteTemplate(w, "aggiungiLibro.html", struct {
		Generi []db.Genere
		Autori []db.Autore
		Values CommonValues
	}{generi, autori, CommonValues{Version}})
}

func HandleGeneraCodici(w http.ResponseWriter, r *http.Request) {

	// Ottiene il cookie
	cookie, err := r.Cookie("access_token")

	// Se non riesce ad ottenerlo reindirizza a /login
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Estrae e controlla il token
	token := []byte(cookie.Value)
	user, err := auth.ParseToken(token)

	// Se l'autenticazione fallisce oppure l'utente non è admin ritorna 401
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if !user.IsAdmin {
		http.Error(w, "You are not an admin!", http.StatusUnauthorized)
		return
	}

	if r.Method == "POST" {
		r.ParseForm()
		codici, ok := r.Form["codici"]
		if !ok || len(codici) < 1 {
			http.Error(w, "Seleziona almeno un elemento!", http.StatusBadRequest)
			return
		}

		var ids []uint32

		for _, codice := range codici {
			id, err := strconv.ParseInt(codice, 2, 32)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			ids = append(ids, uint32(id))
		}

		// Genera la pagina con i codici QR
		page, err := qrcode.GeneratePage(ids)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// Ritorna la pagina
		http.Error(w, page, http.StatusOK)
		return
	}

	libri, err := db.RicercaLibri("", []uint32{}, []uint32{}, -1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	templates.ExecuteTemplate(w, "generaCodici.html", struct {
		Libri  []db.Libro
		Values CommonValues
	}{libri, CommonValues{Version}})
}
