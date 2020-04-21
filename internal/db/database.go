/*
 * database.go
 *
 * Pacchetto per interfacciarsi con il database SQL
 *
 * Copyright (c) 2020 Filippo Casarin <casarin.filippo17@gmail.com>
 *
 * Copyright (c) 2020 Davide Vendramin <davidevendramin5@gmail.com>
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

// Pacchetto per interfacciarsi con il database SQL
package db

import (
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"strings"

	"git.antonionapolitano.eu/napaalm/LilBib/internal/config"
)

//Tabelle del database
type Libro struct {
	codice uint32

	titolo string

	autore uint32

	genere uint32

	prenotato bool

	hashz string
}

type Genere struct {
	codice uint32

	nome string
}

type Autore struct {
	codice uint32

	nome string

	cognome string
}

type Prestito struct {
	codice uint32

	libro uint32

	utente string

	data_prenotazione time.Time

	durata uint32

	data_restituzione time.Time
}

//Funzione per trovare un libro in base al suo codice
func GetLibro(cod uint32) (Libro, error) {

	db, err := sql.Open("mysql", config.Config.SQL.Indirizzo)

	defer db.Close()

	var vuoto Libro //creo un libro vuoto da ritornare in caso di errore

	if err != nil { //se c'è un errore, ritorna un libro vuoto e l'errore

		return vuoto, err

	}

	err = db.Ping() //verifico se il server è ancora disponibile

	if err != nil { //se c'è un errore, ritorna un libro vuoto e l'errore

		return vuoto, err

	}

	q := `SELECT * FROM Libro WHERE codice = ?` //salvo la query che eseguirà l'sql in una variabile stringa

	rows, err := db.Query(q, cod) //applico la query al database. Salvo i risultati in rows

	if err != nil { //se c'è un errore, ritorna un libro vuoto e l'errore

		return vuoto, err

	}

	defer rows.Close() //rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return

	var lib Libro //creo un libro in cui salvare il risultato della query

	for rows.Next() { //rows.Next() scorre tutte le righe trovate dalla query returnando true. Quando le finisce returna false

		if err := rows.Scan(&lib.codice, &lib.titolo, &lib.autore, &lib.genere, &lib.prenotato, &lib.hashz); err != nil { //tramite rows.Scan() salvo i vari risultati nel libro creato in precedenza. In caso di errore ritorno un libro vuoto e l'errore

			return vuoto, err

		}

	}

	if err := rows.Err(); err != nil { //se c'è un errore, ritorna un libro vuoto e l'errore

		return vuoto, err

	}

	return lib, nil //returno il libro trovato e null (null sarebbe l'errore che non è avvenuto)

}

//Funzione per trovare più libri
/*func GetLibri(page uint16, num uint16) ([]Libro, error) {

	db, err := sql.Open("mysql", "buff:rodolfo@/suqi")

	defer db.Close()

	if err != nil { //se c'è un errore, ritorna null e l'errore

		return nil, err

	}

	err = db.Ping() //verifico se il server è ancora disponibile

	if err != nil { //se c'è un errore, ritorna un libro vuoto e l'errore

		return nil, err

	}

	q := `SELECT * FROM Libro LIMIT ?,?` //salvo la query che eseguirà l'sql in una variabile stringa

	rows, err := db.Query(q, (page-1)*num, page*num) //applico la query al database. Salvo i risultati in rows

	if err != nil { //se c'è un errore, ritorna un libro vuoto e l'errore

		return nil, err

	}

	defer rows.Close() //rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return

	var libs []Libro

	for rows.Next() { //rows.Next() scorre tutte le righe trovate dalla query returnando true. Quando le finisce returna false

		var fabrizio Libro //dichiaro variabili temporanee

		if err := rows.Scan(&fabrizio.codice, &fabrizio.titolo, &fabrizio.autore, &fabrizio.genere, &fabrizio.prenotato, &fabrizio.hashz); err != nil { //tramite rows.Scan() salvo i vari risultati nella variabile creata in precedenza. In caso di errore ritorno null e l'errore

			return nil, err

		}

		libs = append(libs, fabrizio) //copio la variabile temporanea nell'ultima posizione dell'array

	}

	if err := rows.Err(); err != nil { //se c'è un errore, ritorna un libro vuoto e l'errore

		return nil, err

	}

	return libs, nil //returno i libri trovati e null (null sarebbe l'errore che non è avvenuto)

}*/

//Funzione per trovare Autori in base all'iniziale del cognome
func GetAutori(iniziale uint8) ([]Autore, error) {

	s := string(iniziale) + "%"

	db, err := sql.Open("mysql", config.Config.SQL.Indirizzo)

	defer db.Close()

	if err != nil { //se c'è un errore, ritorna null e l'errore

		return nil, err

	}

	err = db.Ping() //verifico se il server è ancora disponibile

	if err != nil { //se c'è un errore, ritorna null e l'errore

		return nil, err

	}

	q := `SELECT * FROM Autore WHERE cognome LIKE ?` //salvo la query che eseguirà l'sql in una variabile stringa

	rows, err := db.Query(q, s) //applico la query al database. Salvo i risultati in rows

	if err != nil { //se c'è un errore, ritorna null e l'errore

		return nil, err

	}

	defer rows.Close() //rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return

	var auths []Autore

	for rows.Next() { //rows.Next() scorre tutte le righe trovate dalla query returnando true. Quando le finisce returna false

		var fabrizio Autore //dichiaro variabili temporanee

		if err := rows.Scan(&fabrizio.codice, &fabrizio.nome, &fabrizio.cognome); err != nil { //tramite rows.Scan() salvo i vari risultati nella variabile creata in precedenza. In caso di errore ritorno null e l'errore

			return nil, err

		}

		auths = append(auths, fabrizio) //copio la variabile temporanea nell'ultima posizione dell'array

	}

	if err := rows.Err(); err != nil { //se c'è un errore, ritorna null e l'errore

		return nil, err

	}

	return auths, nil //returno gli autori trovati e null (null sarebbe l'errore che non è avvenuto)

}

//Funzione per trovare tutti i generi esistenti nel database
func GetGeneri() ([]Genere, error) {

	db, err := sql.Open("mysql", config.Config.SQL.Indirizzo)

	defer db.Close()

	if err != nil { //se c'è un errore, ritorna null e l'errore

		return nil, err

	}

	err = db.Ping() //verifico se il server è ancora disponibile

	if err != nil { //se c'è un errore, ritorna null e l'errore

		return nil, err

	}

	q := `SELECT * FROM Genere` //salvo la query che eseguirà l'sql in una variabile stringa

	rows, err := db.Query(q) //applico la query al database. Salvo i risultati in rows

	if err != nil { //se c'è un errore, ritorna null e l'errore

		return nil, err

	}

	defer rows.Close() //rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return

	var gens []Genere

	for rows.Next() { //rows.Next() scorre tutte le righe trovate dalla query returnando true. Quando le finisce returna false

		var fabrizio Genere //dichiaro variabili temporanee

		if err := rows.Scan(&fabrizio.codice, &fabrizio.nome); err != nil { //tramite rows.Scan() salvo i vari risultati nella variabile creata in precedenza. In caso di errore ritorno null e l'errore

			return nil, err

		}

		gens = append(gens, fabrizio) //copio la variabile temporanea nell'ultima posizione dell'array

	}

	if err := rows.Err(); err != nil { //se c'è un errore, ritorna null e l'errore

		return nil, err

	}

	return gens, nil //returno i generi trovati e null (null sarebbe l'errore che non è avvenuto)

}

//Funzione per trovare tutti i prestiti di un utente
func GetPrestiti(utente string) ([]Prestito, error) {

	db, err := sql.Open("mysql", config.Config.SQL.Indirizzo)

	defer db.Close()

	if err != nil { //se c'è un errore, ritorna null e l'errore

		return nil, err

	}

	err = db.Ping() //verifico se il server è ancora disponibile

	if err != nil { //se c'è un errore, ritorna null e l'errore

		return nil, err

	}

	q := `SELECT * FROM Prestito WHERE utente = ?` //salvo la query che eseguirà l'sql in una variabile stringa

	rows, err := db.Query(q, utente) //applico la query al database. Salvo i risultati in rows

	if err != nil { //se c'è un errore, ritorna null e l'errore

		return nil, err

	}

	defer rows.Close() //rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return

	var prests []Prestito

	for rows.Next() { //rows.Next() scorre tutte le righe trovate dalla query returnando true. Quando le finisce returna false

		var data_pre int64 //dichiaro variabili temporanee

		var data_rest int64

		var fabrizio Prestito

		if err := rows.Scan(&fabrizio.codice, &fabrizio.libro, &fabrizio.utente, &data_pre, &fabrizio.durata, &data_rest); err != nil { //tramite rows.Scan() salvo i vari risultati nelle variabili create in precedenza. In caso di errore ritorno null e l'errore

			return nil, err

		}

		fabrizio.data_prenotazione = time.Unix(data_pre, 0) //salvo data_pre in fabrizio convertendola in timestamp unix

		fabrizio.data_restituzione = time.Unix(data_rest, 0) //salvo data_rest in fabrizio convertendola in timestamp unix

		prests = append(prests, fabrizio) //copio la variabile temporanea nell'ultima posizione dell'array

	}

	if err := rows.Err(); err != nil { //se c'è un errore, ritorna null e l'errore

		return nil, err

	}

	return prests, nil //returno i prestiti trovati e null (null sarebbe l'errore che non è avvenuto)

}

func RicercaLibri(nome string, autore, genere []uint32, page, num uint16) ([]Libro, error) {

	db, err := sql.Open("mysql", config.Config.SQL.Indirizzo)

	defer db.Close()

	if err != nil { //se c'è un errore, ritorna null e l'errore

		return nil, err

	}

	err = db.Ping() //verifico se il server è ancora disponibile

	if err != nil { //se c'è un errore, ritorna null e l'errore

		return nil, err

	}

	tags := strings.Split(nome, " ")

	var args []interface{}

	for i := 0; i < len(tags); i++ {

		if len(tags[i]) > 0 {

			tags[i] = "%" + tags[i] + "%"

		}

		args = append(args, tags[i])

	}

	for i := 0; i < len(autore); i++ {

		args = append(args, autore[i])

	}

	for i := 0; i < len(genere); i++ {

		args = append(args, genere[i])

	}

	var q string

	if len(tags) > 0 {

		if len(autore) > 0 {

			if len(genere) > 0 {

				q = `SELECT * FROM Libro WHERE (titolo LIKE ?` + strings.Repeat(` OR titolo LIKE ?`, len(tags)-1) + `) AND (autore = ?` + strings.Repeat(` OR autore = ?`, len(autore)-1) + `) AND (genere = ?` + strings.Repeat(` OR genere = ?`, len(genere)-1) + `) LIMIT ?,?`

			} else {

				q = `SELECT * FROM Libro WHERE (titolo LIKE ?` + strings.Repeat(` OR titolo LIKE ?`, len(tags)-1) + `) AND (autore = ?` + strings.Repeat(` OR autore = ?`, len(autore)-1) + `) LIMIT ?,?`

			}

		} else {

			if len(genere) > 0 {

				q = `SELECT * FROM Libro WHERE (titolo LIKE ?` + strings.Repeat(` OR titolo LIKE ?`, len(tags)-1) + `) AND (genere = ?` + strings.Repeat(` OR genere = ?`, len(genere)-1) + `) LIMIT ?,?`

			} else {

				q = `SELECT * FROM Libro WHERE (titolo LIKE ?` + strings.Repeat(` OR titolo LIKE ?`, len(tags)-1) + `) LIMIT ?,?`

			}

		}

	} else {

		if len(autore) > 0 {

			if len(genere) > 0 {

				q = `SELECT * FROM Libro WHERE (autore = ?` + strings.Repeat(` OR autore = ?`, len(autore)-1) + `) AND (genere = ?` + strings.Repeat(` OR genere = ?`, len(genere)-1) + `) LIMIT ?,?`

			} else {

				q = `SELECT * FROM Libro WHERE (autore = ?` + strings.Repeat(` OR autore = ?`, len(autore)-1) + `) LIMIT ?,?`

			}

		} else {

			if len(genere) > 0 {

				q = `SELECT * FROM Libro WHERE (genere = ?` + strings.Repeat(` OR genere = ?`, len(genere)-1) + `) LIMIT ?,?`

			} else {

				q = `SELECT * FROM Libro LIMIT ?,?`

			}

		}

	}

	args = append(args, (page-1)*num)

	args = append(args, page*num)

	rows, err := db.Query(q, args...)

	if err != nil { //se c'è un errore, ritorna null e l'errore

		return nil, err

	}

	defer rows.Close() //rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return

	var libs []Libro

	for rows.Next() { //rows.Next() scorre tutte le righe trovate dalla query returnando true. Quando le finisce returna false

		var fabrizio Libro //dichiaro variabili temporanee

		if err := rows.Scan(&fabrizio.codice, &fabrizio.titolo, &fabrizio.autore, &fabrizio.genere, &fabrizio.prenotato, &fabrizio.hashz); err != nil { //tramite rows.Scan() salvo i vari risultati nella variabile creata in precedenza. In caso di errore ritorno null e l'errore

			return nil, err

		}

		libs = append(libs, fabrizio) //copio la variabile temporanea nell'ultima posizione dell'array

	}

	if err := rows.Err(); err != nil { //se c'è un errore, ritorna un libro vuoto e l'errore

		return nil, err

	}

	return libs, nil //returno i libri trovati e null (null sarebbe l'errore che non è avvenuto)

}
