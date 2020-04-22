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
	"database/sql"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

//Tabelle del database
type Libro struct {
	Codice    uint32
	Titolo    string
	Autore    uint32
	Genere    uint32
	Prenotato bool
	Hashz     string
}

type Genere struct {
	Codice uint32
	Nome   string
}

type Autore struct {
	Codice  uint32
	Nome    string
	Cognome string
}

type Prestito struct {
	Codice            uint32
	Libro             uint32
	Utente            string
	Data_prenotazione time.Time
	Durata            uint32
	Data_restituzione time.Time
}

//Funzione per trovare un libro in base al suo codice
func GetLibro(cod uint32) (Libro, error) {
	db, err := sql.Open("mysql", config.Config.SQL.Indirizzo)
	//creo un libro vuoto da ritornare in caso di errore
	var vuoto Libro
	//se c'è un errore, ritorna un libro vuoto e l'errore
	if err != nil {
		return vuoto, err
	}
	defer db.Close()

	//verifico se il server è ancora disponibile
	err = db.Ping()
	//se c'è un errore, ritorna un libro vuoto e l'errore
	if err != nil {
		return vuoto, err
	}

	//salvo la query che eseguirà l'sql in una variabile stringa
	q := `SELECT * FROM Libro WHERE codice = ?`
	//applico la query al database. Salvo i risultati in rows
	rows, err := db.Query(q, cod)
	if err != nil { //se c'è un errore, ritorna un libro vuoto e l'errore
		return vuoto, err
	}
	//rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	//creo un libro in cui salvare il risultato della query
	var lib Libro
	//rows.Next() scorre tutte le righe trovate dalla query returnando true. Quando le finisce returna false
	for rows.Next() {
		//tramite rows.Scan() salvo i vari risultati nel libro creato in precedenza. In caso di errore ritorno un libro vuoto e l'errore
		if err := rows.Scan(&lib.Codice, &lib.Titolo, &lib.Autore, &lib.Genere, &lib.Prenotato, &lib.Hashz); err != nil {
			return vuoto, err
		}
	}
	//se c'è un errore, ritorna un libro vuoto e l'errore
	if err := rows.Err(); err != nil {
		return vuoto, err
	}

	//returno il libro trovato e null (null sarebbe l'errore che non è avvenuto)
	return lib, nil
}

//Funzione per trovare più libri
/*func GetLibri(page uint16, num uint16) ([]Libro, error) {
	db, err := sql.Open("mysql", "buff:rodolfo@/suqi")
	if err != nil { //se c'è un errore, ritorna null e l'errore
		return nil, err
	}
	defer db.Close()

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
	db, err := sql.Open("mysql", config.Config.SQL.Indirizzo)
	defer db.Close()
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}

	//verifico se il server è ancora disponibile
	err = db.Ping()
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}

	//creo una stringa s composta dall'iniziale del cognome dell'autore e da %
	s := string(iniziale) + "%"
	//salvo la query che eseguirà l'sql in una variabile stringa
	q := `SELECT * FROM Autore WHERE cognome LIKE ?`
	//applico la query al database. Salvo i risultati in rows
	rows, err := db.Query(q, s)
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}
	//rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var auths []Autore
	//rows.Next() scorre tutte le righe trovate dalla query returnando true. Quando le finisce returna false
	for rows.Next() {
		var fabrizio Autore
		//tramite rows.Scan() salvo i vari risultati nella variabile creata in precedenza. In caso di errore ritorno null e l'errore
		if err := rows.Scan(&fabrizio.Codice, &fabrizio.Nome, &fabrizio.Cognome); err != nil {
			return nil, err
		}
		//copio la variabile temporanea nell'ultima posizione dell'array
		auths = append(auths, fabrizio)
	}
	//se c'è un errore, ritorna null e l'errore
	if err := rows.Err(); err != nil {
		return nil, err
	}

	//returno gli autori trovati e null (null sarebbe l'errore che non è avvenuto)
	return auths, nil
}

//Funzione per trovare tutti i generi esistenti nel database
func GetGeneri() ([]Genere, error) {
	db, err := sql.Open("mysql", config.Config.SQL.Indirizzo)
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}
	defer db.Close()

	//verifico se il server è ancora disponibile
	err = db.Ping()
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}

	//salvo la query che eseguirà l'sql in una variabile stringa
	q := `SELECT * FROM Genere`
	//applico la query al database. Salvo i risultati in rows
	rows, err := db.Query(q)
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}
	//rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var gens []Genere
	//rows.Next() scorre tutte le righe trovate dalla query returnando true. Quando le finisce returna false
	for rows.Next() {
		var fabrizio Genere
		//tramite rows.Scan() salvo i vari risultati nella variabile creata in precedenza. In caso di errore ritorno null e l'errore
		if err := rows.Scan(&fabrizio.Codice, &fabrizio.Nome); err != nil {
			return nil, err
		}
		//copio la variabile temporanea nell'ultima posizione dell'array
		gens = append(gens, fabrizio)
	}
	//se c'è un errore, ritorna null e l'errore
	if err := rows.Err(); err != nil {
		return nil, err
	}

	//returno i generi trovati e null (null sarebbe l'errore che non è avvenuto)
	return gens, nil
}

//Funzione per trovare tutti i prestiti di un utente
func GetPrestiti(utente string) ([]Prestito, error) {
	db, err := sql.Open("mysql", config.Config.SQL.Indirizzo)
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}
	defer db.Close()

	//verifico se il server è ancora disponibile
	err = db.Ping()
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}

	//salvo la query che eseguirà l'sql in una variabile stringa
	q := `SELECT * FROM Prestito WHERE utente = ?`
	//applico la query al database. Salvo i risultati in rows
	rows, err := db.Query(q, utente)
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}
	//rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var prests []Prestito
	//rows.Next() scorre tutte le righe trovate dalla query returnando true. Quando le finisce returna false
	for rows.Next() {
		var data_pre int64
		var data_rest int64
		var fabrizio Prestito
		//tramite rows.Scan() salvo i vari risultati nelle variabili create in precedenza. In caso di errore ritorno null e l'errore
		if err := rows.Scan(&fabrizio.Codice, &fabrizio.Libro, &fabrizio.Utente, &data_pre, &fabrizio.Durata, &data_rest); err != nil {
			return nil, err
		}
		//salvo data_pre in fabrizio convertendola in timestamp unix
		fabrizio.Data_prenotazione = time.Unix(data_pre, 0)
		//salvo data_rest in fabrizio convertendola in timestamp unix
		fabrizio.Data_restituzione = time.Unix(data_rest, 0)
		//copio la variabile temporanea nell'ultima posizione dell'array
		prests = append(prests, fabrizio)
	}
	//se c'è un errore, ritorna null e l'errore
	if err := rows.Err(); err != nil {
		return nil, err
	}

	//returno i prestiti trovati e null (null sarebbe l'errore che non è avvenuto)
	return prests, nil
}

func RicercaLibri(nome string, autore, genere []uint32, page, num uint16) ([]Libro, error) {
	db, err := sql.Open("mysql", config.Config.SQL.Indirizzo)
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}
	defer db.Close()
	//verifico se il server è ancora disponibile
	err = db.Ping()
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}

	//divido la stringa nome in vari tag e poi aggiungo i vari argomenti alla slice "args"
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

	//esamino tutti i casi possibili di richiesta, scegliendo la query giusta per ogni situazione possibile
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
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}
	//rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var libs []Libro
	//rows.Next() scorre tutte le righe trovate dalla query returnando true. Quando le finisce returna false
	for rows.Next() {
		var fabrizio Libro
		//tramite rows.Scan() salvo i vari risultati nella variabile creata in precedenza. In caso di errore ritorno null e l'errore
		if err := rows.Scan(&fabrizio.Codice, &fabrizio.Titolo, &fabrizio.Autore, &fabrizio.Genere, &fabrizio.Prenotato, &fabrizio.Hashz); err != nil {
			return nil, err
		}
		//copio la variabile temporanea nell'ultima posizione dell'array
		libs = append(libs, fabrizio)
	}
	//se c'è un errore, ritorna null e l'errore
	if err := rows.Err(); err != nil {
		return nil, err
	}

	//returno i libri trovati e null (null sarebbe l'errore che non è avvenuto)
	return libs, nil
}

func RicercaAutori(nome string) ([]Autore, error) {
	db, err := sql.Open("mysql", config.Config.SQL.Indirizzo)
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}
	defer db.Close()

	//verifico se il server è ancora disponibile
	err = db.Ping()
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}

	//divido la stringa nome in vari tag e poi li aggiungo alla slice "args"
	tags := strings.Split(nome, " ")
	var args []interface{}
	for j := 0; j < 2; j++ {
		for i := 0; i < len(tags); i++ {
			if len(tags[i]) > 0 {
				//i % servono per dire all'SQL di cercare la stringa in qualsiasi posizione
				args = append(args, "%"+tags[i]+"%")
			}
		}
	}

	q := `SELECT * FROM Autore WHERE nome LIKE ?` + strings.Repeat(` OR nome LIKE ?`, len(tags)-1) + strings.Repeat(` OR cognome LIKE ?`, len(tags))
	rows, err := db.Query(q, args...)
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}
	//rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var auths []Autore
	for rows.Next() {
		var fabrizio Autore
		//tramite rows.Scan() salvo i vari risultati nella variabile creata in precedenza. In caso di errore ritorno null e l'errore
		if err := rows.Scan(&fabrizio.Codice, &fabrizio.Nome, &fabrizio.Cognome); err != nil {
			return nil, err
		}
		//copio la variabile temporanea nell'ultima posizione dell'array
		auths = append(auths, fabrizio)
	}

	//se c'è un errore, ritorna null e l'errore
	if err := rows.Err(); err != nil {
		return nil, err
	}

	//returno gli autori trovati e null (null sarebbe l'errore che non è avvenuto)
	return auths, nil
}

func RicercaGeneri(nome string) ([]Genere, error) {
	db, err := sql.Open("mysql", config.Config.SQL.Indirizzo)
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}
	defer db.Close()

	//verifico se il server è ancora disponibile
	err = db.Ping()
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}

	//divido la stringa nome in vari tag e poi li aggiungo alla slice "args"
	tags := strings.Split(nome, " ")
	var args []interface{}
	for i := 0; i < len(tags); i++ {
		if len(tags[i]) > 0 {
			//i % servono per dire all'SQL di cercare la stringa in qualsiasi posizione
			args = append(args, "%"+tags[i]+"%")
		}
	}

	q := `SELECT * FROM Genere WHERE nome LIKE ?` + strings.Repeat(` OR nome LIKE ?`, len(tags)-1)
	rows, err := db.Query(q, args...)
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}
	//rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var gens []Genere
	for rows.Next() {
		var fabrizio Genere
		//tramite rows.Scan() salvo i vari risultati nella variabile creata in precedenza. In caso di errore ritorno null e l'errore
		if err := rows.Scan(&fabrizio.Codice, &fabrizio.Nome); err != nil {
			return nil, err
		}
		//copio la variabile temporanea nell'ultima posizione dell'array
		gens = append(gens, fabrizio)
	}

	//se c'è un errore, ritorna null e l'errore
	if err := rows.Err(); err != nil {
		return nil, err
	}

	//returno i generi trovati e null (null sarebbe l'errore che non è avvenuto)
	return gens, nil
}

func AddLibro(titolo string, autore, genere uint32) (uint32, error) {
	db, err := sql.Open("mysql", config.Config.SQL.Indirizzo)
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return 0, err
	}
	defer db.Close()

	//verifico se il server è ancora disponibile
	err = db.Ping()
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return 0, err
	}

	q := `INSERT INTO Libro VALUES (null, ?, ?, ?, false, null);
		  SELECT LAST_INSERT_ID();`
	rows, err := db.Query(q, titolo, autore, genere)
	//se c'è un errore, ritorna null e l'errore
	if err != nil {
		return 0, err
	}
	//rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var id uint32
	rows.Next()
	//tramite rows.Scan() salvo il risultato nella variabile creata in precedenza. In caso di errore ritorno null e l'errore
	if err := rows.Scan(&id); err != nil {
		return 0, err
	}
	//se c'è un errore, ritorna null e l'errore
	if err := rows.Err(); err != nil {
		return 0, err
	}

	//returno l'id del libro inserito e null (null sarebbe l'errore che non è avvenuto)
	return id, nil
}
