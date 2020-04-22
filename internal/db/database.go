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
 * <http://Www.gnu.org/licenses/>.
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

var db_Connection *sql.DB

//Funzione per inizializzare il database
func InizializzaDB() error {
	var err error
	db_Connection, err = sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", config.Config.SQL.Username, config.Config.SQL.Password, config.Config.SQL.Indirizzo, config.Config.SQL.Database))
	return err
}

//Funzione per chiudere il database
func ChiudiDB() {
	db_Connection.Close()
}

//Funzione per trovare un libro in base al suo codice
func GetLibro(cod uint32) (Libro, error) {
	//Verifico se il server è ancora disponibile
	//Se c'è un errore, ritorna un libro vuoto e l'errore
	if err := db_Connection.Ping(); err != nil {
		return Libro{}, err
	}

	//Salvo la query che eseguirà l'sql in una variabile stringa
	q := `SELECT * FROM Libro WHERE codice = ?`
	//Applico la query al database. Salvo i risultati in rows
	rows, err := db_Connection.Query(q, cod)
	//Se c'è un errore, ritorna un libro vuoto e l'errore
	if err != nil {
		return Libro{}, err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	//Creo un libro in cui salvare il risultato della query
	var lib Libro
	//Rows.Next() scorre tutte le righe trovate dalla query returnando true. Quando le finisce returna false
	for rows.Next() {
		//Tramite rows.Scan() salvo i vari risultati nel libro creato in precedenza. In caso di errore ritorno un libro vuoto e l'errore
		if err := rows.Scan(&lib.Codice, &lib.Titolo, &lib.Autore, &lib.Genere, &lib.Prenotato, &lib.Hashz); err != nil {
			return Libro{}, err
		}
	}
	//Se c'è un errore, ritorna un libro vuoto e l'errore
	if err := rows.Err(); err != nil {
		return Libro{}, err
	}

	//Returno il libro trovato e null (null sarebbe l'errore che non è avvenuto)
	return lib, nil
}

//Funzione per trovare più libri
/*func GetLibri(page uint16, num uint16) ([]Libro, error) {
	//Verifico se il server è ancora disponibile
	err := db_Connection.Ping()
	//Se c'è un errore, ritorna un libro vuoto e l'errore
	if err != nil {
		return nil, err
	}

	//Salvo la query che eseguirà l'sql in una variabile stringa
	q := `SELECT * FROM Libro LIMIT ?,?`
	//Applico la query al database. Salvo i risultati in rows
	rows, err := db_Connection.Query(q, (page-1)*num, page*num)
	//Se c'è un errore, ritorna un libro vuoto e l'errore
	if err != nil {
		return nil, err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var libs []Libro
	//Rows.Next() scorre tutte le righe trovate dalla query returnando true. Quando le finisce returna false
	for rows.Next() {
		//Dichiaro variabili temporanee
		var fabrizio Libro
		//Tramite rows.Scan() salvo i vari risultati nella variabile creata in precedenza. In caso di errore ritorno null e l'errore
		if err := rows.Scan(&fabrizio.codice, &fabrizio.titolo, &fabrizio.autore, &fabrizio.genere, &fabrizio.prenotato, &fabrizio.hashz); err != nil {
			return nil, err
		}
		//Copio la variabile temporanea nell'ultima posizione dell'array
		libs = append(libs, fabrizio)
	}
	//Se c'è un errore, ritorna un libro vuoto e l'errore
	if err := rows.Err(); err != nil {
		return nil, err
	}
	//Returno i libri trovati e null (null sarebbe l'errore che non è avvenuto)
	return libs, nil
}*/

//Funzione per trovare Autori in base all'iniziale del cognome
func GetAutori(iniziale uint8) ([]Autore, error) {
	//Verifico se il server è ancora disponibile
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
		return nil, err
	}

	//Creo una stringa s composta dall'iniziale del cognome dell'autore e da %
	s := string(iniziale) + "%"
	//Salvo la query che eseguirà l'sql in una variabile stringa
	q := `SELECT * FROM Autore WHERE cognome LIKE ?`
	//Applico la query al database. Salvo i risultati in rows
	rows, err := db_Connection.Query(q, s)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var auths []Autore
	//Rows.Next() scorre tutte le righe trovate dalla query returnando true. Quando le finisce returna false
	for rows.Next() {
		var fabrizio Autore
		//Tramite rows.Scan() salvo i vari risultati nella variabile creata in precedenza. In caso di errore ritorno null e l'errore
		if err := rows.Scan(&fabrizio.Codice, &fabrizio.Nome, &fabrizio.Cognome); err != nil {
			return nil, err
		}
		//Copio la variabile temporanea nell'ultima posizione dell'array
		auths = append(auths, fabrizio)
	}
	//Se c'è un errore, ritorna null e l'errore
	if err := rows.Err(); err != nil {
		return nil, err
	}

	//Returno gli autori trovati e null (null sarebbe l'errore che non è avvenuto)
	return auths, nil
}

//Funzione per trovare tutti i generi esistenti nel database
func GetGeneri() ([]Genere, error) {
	//Verifico se il server è ancora disponibile
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
		return nil, err
	}

	//Salvo la query che eseguirà l'sql in una variabile stringa
	q := `SELECT * FROM Genere`
	//Applico la query al database. Salvo i risultati in rows
	rows, err := db_Connection.Query(q)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var gens []Genere
	//Rows.Next() scorre tutte le righe trovate dalla query returnando true. Quando le finisce returna false
	for rows.Next() {
		var fabrizio Genere
		//Tramite rows.Scan() salvo i vari risultati nella variabile creata in precedenza. In caso di errore ritorno null e l'errore
		if err := rows.Scan(&fabrizio.Codice, &fabrizio.Nome); err != nil {
			return nil, err
		}
		//Copio la variabile temporanea nell'ultima posizione dell'array
		gens = append(gens, fabrizio)
	}
	//Se c'è un errore, ritorna null e l'errore
	if err := rows.Err(); err != nil {
		return nil, err
	}

	//Returno i generi trovati e null (null sarebbe l'errore che non è avvenuto)
	return gens, nil
}

//Funzione per trovare tutti i prestiti di un utente
func GetPrestiti(utente string) ([]Prestito, error) {
	//Verifico se il server è ancora disponibile
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
		return nil, err
	}

	//Salvo la query che eseguirà l'sql in una variabile stringa
	q := `SELECT * FROM Prestito WHERE utente = ?`
	//Applico la query al database. Salvo i risultati in rows
	rows, err := db_Connection.Query(q, utente)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var prests []Prestito
	//Rows.Next() scorre tutte le righe trovate dalla query returnando true. Quando le finisce returna false
	for rows.Next() {
		var data_pre int64
		var data_rest int64
		var fabrizio Prestito
		//Tramite rows.Scan() salvo i vari risultati nelle variabili create in precedenza. In caso di errore ritorno null e l'errore
		if err := rows.Scan(&fabrizio.Codice, &fabrizio.Libro, &fabrizio.Utente, &data_pre, &fabrizio.Durata, &data_rest); err != nil {
			return nil, err
		}
		//Salvo data_pre in fabrizio convertendola in timestamp unix
		fabrizio.Data_prenotazione = time.Unix(data_pre, 0)
		//Salvo data_rest in fabrizio convertendola in timestamp unix
		fabrizio.Data_restituzione = time.Unix(data_rest, 0)
		//Copio la variabile temporanea nell'ultima posizione dell'array
		prests = append(prests, fabrizio)
	}
	//Se c'è un errore, ritorna null e l'errore
	if err := rows.Err(); err != nil {
		return nil, err
	}

	//Returno i prestiti trovati e null (null sarebbe l'errore che non è avvenuto)
	return prests, nil
}

//Funzione per la ricerca dei libri
func RicercaLibri(nome string, autore, genere []uint32, page, num uint16) ([]Libro, error) {
	//Verifico se il server è ancora disponibile
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
		return nil, err
	}

	//Divido la stringa nome in vari tag e poi aggiungo i vari argomenti alla slice "args"
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

	//Esamino tutti i casi possibili di richiesta, scegliendo la query giusta per ogni situazione possibile
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
	args = append(args, page*num, (page+1)*num)

	rows, err := db_Connection.Query(q, args...)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var libs []Libro
	//Rows.Next() scorre tutte le righe trovate dalla query returnando true. Quando le finisce returna false
	for rows.Next() {
		var fabrizio Libro
		//Tramite rows.Scan() salvo i vari risultati nella variabile creata in precedenza. In caso di errore ritorno null e l'errore
		if err := rows.Scan(&fabrizio.Codice, &fabrizio.Titolo, &fabrizio.Autore, &fabrizio.Genere, &fabrizio.Prenotato, &fabrizio.Hashz); err != nil {
			return nil, err
		}
		//Copio la variabile temporanea nell'ultima posizione dell'array
		libs = append(libs, fabrizio)
	}
	//Se c'è un errore, ritorna null e l'errore
	if err := rows.Err(); err != nil {
		return nil, err
	}

	//Returno i libri trovati e null (null sarebbe l'errore che non è avvenuto)
	return libs, nil
}

//Funzione per la ricerca degli autori
func RicercaAutori(nome string) ([]Autore, error) {
	//Verifico se il server è ancora disponibile
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
		return nil, err
	}

	//Divido la stringa nome in vari tag e poi li aggiungo alla slice "args"
	tags := strings.Split(nome, " ")
	var args []interface{}
	for j := uint8(0); j < 2; j++ {
		for _, tag := range tags {
			if len(tag) > 0 {
				//I % servono per dire all'SQL di cercare la stringa in qualsiasi posizione
				args = append(args, "%"+tag+"%")
			}
		}
	}

	q := `SELECT * FROM Autore WHERE nome LIKE ?` + strings.Repeat(` OR nome LIKE ?`, len(tags)-1) + strings.Repeat(` OR cognome LIKE ?`, len(tags))
	rows, err := db_Connection.Query(q, args...)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var auths []Autore
	for rows.Next() {
		var fabrizio Autore
		//Tramite rows.Scan() salvo i vari risultati nella variabile creata in precedenza. In caso di errore ritorno null e l'errore
		if err := rows.Scan(&fabrizio.Codice, &fabrizio.Nome, &fabrizio.Cognome); err != nil {
			return nil, err
		}
		//Copio la variabile temporanea nell'ultima posizione dell'array
		auths = append(auths, fabrizio)
	}

	//Se c'è un errore, ritorna null e l'errore
	if err := rows.Err(); err != nil {
		return nil, err
	}

	//Returno gli autori trovati e null (null sarebbe l'errore che non è avvenuto)
	return auths, nil
}

//Funzione per la ricerca dei generi
func RicercaGeneri(nome string) ([]Genere, error) {
	//Verifico se il server è ancora disponibile
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
		return nil, err
	}

	//Divido la stringa nome in vari tag e poi li aggiungo alla slice "args"
	tags := strings.Split(nome, " ")
	var args []interface{}
	for _, tag := range rags {
		if len(tag) > 0 {
			//I % servono per dire all'SQL di cercare la stringa in qualsiasi posizione
			args = append(args, "%"+tag+"%")
		}
	}

	q := `SELECT * FROM Genere WHERE nome LIKE ?` + strings.Repeat(` OR nome LIKE ?`, len(tags)-1)
	rows, err := db_Connection.Query(q, args...)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return nil, err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var gens []Genere
	for rows.Next() {
		var fabrizio Genere
		//Tramite rows.Scan() salvo i vari risultati nella variabile creata in precedenza. In caso di errore ritorno null e l'errore
		if err := rows.Scan(&fabrizio.Codice, &fabrizio.Nome); err != nil {
			return nil, err
		}
		//Copio la variabile temporanea nell'ultima posizione dell'array
		gens = append(gens, fabrizio)
	}

	//Se c'è un errore, ritorna null e l'errore
	if err := rows.Err(); err != nil {
		return nil, err
	}

	//Returno i generi trovati e null (null sarebbe l'errore che non è avvenuto)
	return gens, nil
}

//Funzione per aggiungere un libro

func AddLibro(titolo string, autore, genere uint32) (uint32, error) {
	//Verifico se il server è ancora disponibile
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
		return 0, err
	}

	//Preparo il database per la query
	stmt, err := db_Connection.Prepare(`INSERT INTO Libro VALUES (null, ?, ?, ?, false, "")`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	//Eseguo la query e ne salvo i risultati
	res, err := stmt.Exec(titolo, autore, genere)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return 0, err
	}

	//Trovo l'id del libro appena inserito
	id, err := res.LastInsertId()
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return 0, err
	}

	//Returno l'id del libro inserito e null (null sarebbe l'errore che non è avvenuto)
	return uint32(id), nil
}

//Funzione per aggiungere un genere
func AddGenere(nome string) (uint32, error) {
	//Verifico se il server è ancora disponibile
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
		return 0, err
	}

	//Preparo il database per la query
	stmt, err := db_Connection.Prepare(`INSERT INTO Genere VALUES (null, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	//Eseguo la query e ne salvo i risultati
	res, err := stmt.Exec(nome)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return 0, err
	}

	//Trovo l'id del genere appena inserito
	id, err := res.LastInsertId()
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return 0, err
	}

	//Returno l'id del libro inserito e null (null sarebbe l'errore che non è avvenuto)
	return uint32(id), nil
}

//Funzione per aggiungere un autore
func AddAutore(nome, cognome string) (uint32, error) {
	//Verifico se il server è ancora disponibile
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
		return 0, err
	}

	//Preparo il database per la query
	stmt, err := db_Connection.Prepare(`INSERT INTO Autore VALUES (null, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	//Eseguo la query e ne salvo i risultati
	res, err := stmt.Exec(nome, cognome)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return 0, err
	}

	//Trovo l'id del genere appena inserito
	id, err := res.LastInsertId()
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return 0, err
	}

	//Returno l'id del libro inserito e null (null sarebbe l'errore che non è avvenuto)
	return uint32(id), nil
}

func AddPrestito(libro uint32, utente string, data_prenotazione time.Time, durata uint32) (uint32, error) {
	//Verifico se il server è ancora disponibile
	err := db_Connection.Ping()
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return 0, err
	}

	//Preparo il database per la query
	stmt, err := db_Connection.Prepare(`INSERT INTO Prestito VALUES (null, ?, ?, ?, ?, null)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	//Eseguo la query e ne salvo i risultati
	res, err := stmt.Exec(libro, utente, data_prenotazione.Unix(), durata)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return 0, err
	}

	//Trovo l'id del genere appena inserito
	id, err := res.LastInsertId()
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return 0, err
	}

	//Returno l'id del libro inserito e null (null sarebbe l'errore che non è avvenuto)
	return uint32(id), nil
}

//Funzione per impostare l'hash di un libro
func SetHash(codice uint32, hash string) error {
	//Verifico se il server è ancora disponibile
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
		return err
	}

	q := `UPDATE Libro
		  SET hashz = ?
		  WHERE codice = ?`
	rows, err := db_Connection.Query(q, hash, codice)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	return nil
}

//Funzione per impostare la restituzione
func SetRestituzione(prestito uint32, data_restituzione time.Time) error {
	//Verifico se il server è ancora disponibile
	err := db_Connection.Ping()
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return err
	}

	q := `UPDATE Prestito
		  SET data_restituzione = ?
		  WHERE codice = ?`
	rows, err := db_Connection.Query(q, data_restituzione.Unix(), prestito)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	return nil
}

//Funzione per rimuovere un libro
func RemoveLibro(codice uint32) error {
	//Verifico se il server è ancora disponibile
	err := db_Connection.Ping()
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return err
	}

	q := `DELETE FROM Libro
		  WHERE codice = ?`
	rows, err := db_Connection.Query(q, codice)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	return nil
}

//Funzione per rimuovere un genere
func RemoveGenere(codice uint32) error {
	//Verifico se il server è ancora disponibile
	err := db_Connection.Ping()
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return err
	}

	q := `DELETE FROM Genere
		  WHERE codice = ?`
	rows, err := db_Connection.Query(q, codice)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	return nil
}

//Funzione per rimuovere un autore
func RemoveAutore(codice uint32) error {
	//Verifico se il server è ancora disponibile
	err := db_Connection.Ping()
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return err
	}

	q := `DELETE FROM Autore
		  WHERE codice = ?`
	rows, err := db_Connection.Query(q, codice)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	return nil
}

//Funzione per rimuovere un prestito
func RemovePrestito(codice uint32) error {
	//Verifico se il server è ancora disponibile
	err := db_Connection.Ping()
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return err
	}

	q := `DELETE FROM Prestito
		  WHERE codice = ?`
	rows, err := db_Connection.Query(q, codice)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	return nil
}
