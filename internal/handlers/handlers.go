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
	"fmt"
	"git.napaalm.xyz/napaalm/LilBib/internal/auth"
	"git.napaalm.xyz/napaalm/LilBib/internal/config"
	"git.napaalm.xyz/napaalm/LilBib/internal/db"
	"git.napaalm.xyz/napaalm/LilBib/internal/hash"
	"git.napaalm.xyz/napaalm/LilBib/internal/qrcode"
	"net/http"
	"net/url"
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
	prestito, err := db.GetCurrentPrestito(idLibro)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	prestiti, err := db.GetPrestitiLibro(idLibro)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	templates.ExecuteTemplate(w, "libro.html", struct {
		Libro           db.Libro
		Prestiti        []db.Prestito
		CurrentPrestito db.NullPrestito
		Values          CommonValues
	}{libro, prestiti, prestito, CommonValues{Version}})
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
	page := uint16(pageParsed)

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

	libriTot, err := db.RicercaLibriNoPage(titolo, idsAutori, idsGeneri)
	if float32(page) > float32(len(libriTot))/float32(config.Config.Generale.LunghezzaPagina) {
		http.Redirect(w, r, fmt.Sprintf("/libri/%d?titolo=%s&autori=%s&generi=%s", len(libriTot)/int(config.Config.Generale.LunghezzaPagina), titolo, nomeAutore, nomeGenere), http.StatusSeeOther)
		return
	}

	libri, err := db.RicercaLibri(titolo, idsAutori, idsGeneri, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if page == 0 {
		if float32(page+1) >= (float32(len(libriTot)) / float32(config.Config.Generale.LunghezzaPagina)) {
			templates.ExecuteTemplate(w, "libri.html", struct {
				PaginaPrec uint16
				Pagina     uint16
				PaginaSucc uint16
				Titolo     string
				Autori     string
				Generi     string
				Libri      []db.Libro
				Values     CommonValues
			}{page, page + 1, page, titolo, nomeAutore, nomeGenere, libri, CommonValues{Version}})
		} else {
			templates.ExecuteTemplate(w, "libri.html", struct {
				PaginaPrec uint16
				Pagina     uint16
				PaginaSucc uint16
				Titolo     string
				Autori     string
				Generi     string
				Libri      []db.Libro
				Values     CommonValues
			}{page, page + 1, page + 1, titolo, nomeAutore, nomeGenere, libri, CommonValues{Version}})

		}
	} else {
		if float32(page+1) >= (float32(len(libri)) / float32(config.Config.Generale.LunghezzaPagina)) {
			templates.ExecuteTemplate(w, "libri.html", struct {
				PaginaPrec uint16
				Pagina     uint16
				PaginaSucc uint16
				Titolo     string
				Autori     string
				Generi     string
				Libri      []db.Libro
				Values     CommonValues
			}{page - 1, page + 1, page, titolo, nomeAutore, nomeGenere, libri, CommonValues{Version}})
		} else {
			templates.ExecuteTemplate(w, "libri.html", struct {
				PaginaPrec uint16
				Pagina     uint16
				PaginaSucc uint16
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
	// Usa o meno un provider SSO
	sso := config.Config.Autenticazione.SSO

	if r.Method == "POST" && !sso {
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

		// Ottiene la configurazione per i cookie
		fqdn := config.Config.Generale.FQDN
		secure := config.Config.Autenticazione.SecureCookies

		// Crea e imposta il cookie
		cookie := http.Cookie{
			Name:   "access_token",
			Value:  string(token),
			Domain: fqdn,
			MaxAge: 86400, // 24 ore
			Secure: secure,
		}
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

	if sso {
		// Imposta dominio ed eventualmente porta
		nextURL := config.Config.Generale.FQDN
		if nextURL == "localhost" {
			nextURL = nextURL + config.Config.Generale.Porta
		}

		// Imposta il percorso
		nextURL = nextURL + "/utente"

		// Imposta lo schema
		if config.Config.Autenticazione.SecureCookies {
			nextURL = "https://" + nextURL
		} else {
			nextURL = "http://" + nextURL
		}

		// Crea la query
		query := "?next=" + url.QueryEscape(nextURL)

		// Genera la risposta
		http.Redirect(w, r, config.Config.Autenticazione.SSOURL+query, http.StatusSeeOther)
	} else {
		templates.ExecuteTemplate(w, "login.html", struct {
			Values CommonValues
		}{CommonValues{Version}})
	}
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
		deleteCookie(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	prestiti, err := db.GetPrestitiUtente(utente.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prestitiTitoli := make([]PrestitoTitolo, len(prestiti))
	for index, prestito := range prestiti {
		libro, err := db.GetLibro(prestito.Libro)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		prestitiTitoli[index].Prestito = prestiti[index]
		prestitiTitoli[index].Titolo = libro.Titolo
	}

	isAdmin := utente.Username == config.Config.Generale.AdminUser
	templates.ExecuteTemplate(w, "utente.html", struct {
		Utente         auth.UserInfo
		IsAdmin        bool
		PrestitiTitoli []PrestitoTitolo
		Values         CommonValues
	}{utente, isAdmin, prestitiTitoli, CommonValues{Version}})
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
	_, err = db.GetCurrentPrestito(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Imposta la restituzione
	err = db.SetRestituzione(id)

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
	} else if user.Username != config.Config.Generale.AdminUser {
		http.Error(w, "You are not an admin!", http.StatusUnauthorized)
		return
	}

	if r.Method == "POST" {
		r.ParseForm()
		tipo_list, ok := r.Form["tipo"]
		if !ok {
			http.Error(w, "tipo non definito", http.StatusBadRequest)
			return
		}

		tipo := tipo_list[0]

		if tipo == "genere" {
			genere, ok := r.Form["genere"]
			if !ok {
				http.Error(w, "genere non definito", http.StatusBadRequest)
				return
			}

			if _, err := db.AddGenere(genere[0]); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/admin/aggiungiLibro", http.StatusSeeOther)
			return
		}

		if tipo == "autore" {
			nome, ok0 := r.Form["nome"]
			cognome, ok1 := r.Form["cognome"]
			if !ok0 || !ok1 {
				http.Error(w, "autore non definito", http.StatusBadRequest)
				return
			}

			if _, err := db.AddAutore(nome[0], cognome[0]); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/admin/aggiungiLibro", http.StatusSeeOther)
			return
		}

		if tipo == "libro" {
			genere_str, ok0 := r.Form["genere"]
			autore_str, ok1 := r.Form["autore"]
			libro, ok2 := r.Form["nome"]

			if !ok0 || !ok1 || !ok2 {
				http.Error(w, "valori non definiti", http.StatusBadRequest)
				return
			}

			genere, err := strconv.ParseUint(genere_str[0], 10, 32)
			if err != nil {
				http.Error(w, "genere non è un intero", http.StatusBadRequest)
				return
			}

			autore, err := strconv.ParseUint(autore_str[0], 10, 32)
			if err != nil {
				http.Error(w, "autore non è un intero", http.StatusBadRequest)
				return
			}

			if _, err := db.AddLibro(libro[0], uint32(autore), uint32(genere)); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/admin/aggiungiLibro", http.StatusSeeOther)
			return
		}

		http.Error(w, "valore invalido per tipo", http.StatusBadRequest)
		return
	}

	generi, err := db.GetGeneri()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	autori, err := db.GetTuttiAutori()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	} else if user.Username != config.Config.Generale.AdminUser {
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
			id, err := strconv.ParseInt(codice, 10, 32)
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
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(page))
		return
	}

	libri, err := db.GetLibri()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	templates.ExecuteTemplate(w, "generaCodici.html", struct {
		Libri  []db.Libro
		Values CommonValues
	}{libri, CommonValues{Version}})
}

// Percorso: /logout
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	deleteCookie(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Elimina il cookie con il token di accesso
func deleteCookie(w http.ResponseWriter, r *http.Request) {
	// Usa o meno un provider SSO
	sso := config.Config.Autenticazione.SSO

	// Delega il logout al provider SSO oppure elimina direttamente il cookie
	if sso {
		// Imposta dominio ed eventualmente porta
		nextURL := config.Config.Generale.FQDN
		if nextURL == "localhost" {
			nextURL = nextURL + config.Config.Generale.Porta
		}

		// Imposta lo schema
		if config.Config.Autenticazione.SecureCookies {
			nextURL = "https://" + nextURL
		} else {
			nextURL = "http://" + nextURL
		}

		// Crea la query
		query := "?next=" + url.QueryEscape(nextURL)

		// Genera la risposta
		http.Redirect(w, r, config.Config.Autenticazione.SSOURL+"/logout"+query, http.StatusSeeOther)
	} else {
		// Ottiene la configurazione per i cookie
		fqdn := config.Config.Generale.FQDN
		secure := config.Config.Autenticazione.SecureCookies

		// Crea e imposta il cookie
		cookie := http.Cookie{
			Name:   "access_token",
			Value:  "",
			Domain: fqdn,
			MaxAge: -1,
			Secure: secure,
		}
		http.SetCookie(w, &cookie)
	}
}
