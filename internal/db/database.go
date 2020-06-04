/*
 * database.go
 *
 * Pacchetto per interfacciarsi con il database SQL
 *
 * Copyright (c) 2020 Filippo Casarin <casarin.filippo17@gmail.com>
 * Copyright (c) 2020 Davide Vendramin <natalianatiche02@gmail.com>
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
	"fmt"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

//Tabelle del database
type Libro struct {
	Codice        uint32
	Titolo        string
	NomeAutore    string
	CognomeAutore string
	Genere        string
	Prenotato     bool
	Hashz         string
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

type NoCurrentPrestitoError struct {
	codice uint32
}

func (e *NoCurrentPrestitoError) Error() string {
	return fmt.Sprintf("No current Prestito for libro with codice %d.", e.codice)
}

var db_Connection *sql.DB

//Funzione per inizializzare il database
func InizializzaDB() (err error) {
	db_Connection, err = sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", config.Config.SQL.Username, config.Config.SQL.Password, config.Config.SQL.Indirizzo, config.Config.SQL.Database))
	return
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
	q := `SELECT Libro.Codice,Libro.Titolo,Autore.Nome,Autore.Cognome,Genere.Nome,Prenotato,Hashz FROM Libro,Autore,Genere WHERE Libro.Autore = Autore.Codice AND Libro.Genere = Genere.Codice AND Libro.Codice = ?`
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
		if err := rows.Scan(&lib.Codice, &lib.Titolo, &lib.NomeAutore, &lib.CognomeAutore, &lib.Genere, &lib.Prenotato, &lib.Hashz); err != nil {
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

//Funzione per trovare l'assegnatario di un libro
func GetAssegnatario(cod uint32) (string, error) {
	//Verifico se il server è ancora disponibile
	//Se c'è un errore, ritorna un libro vuoto e l'errore
	if err := db_Connection.Ping(); err != nil {
		return "", err
	}

	//Salvo la query che eseguirà l'sql in una variabile stringa
	qpre := `SELECT COUNT(*) FROM Prestito WHERE Prestito.Libro = ?`
	//Applico la query al database. Salvo i risultati in rows
	rows, err := db_Connection.Query(qpre, cod)
	//Se c'è un errore, ritorna un libro vuoto e l'errore
	if err != nil {
		return "", err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var count uint32
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return "", err
		}
	}
	if count == 0 {
		return "", nil
	}

	//Salvo la query che eseguirà l'sql in una variabile stringa
	q := `SELECT Prestito.Utente FROM Prestito WHERE Prestito.Libro = ?`
	//Applico la query al database. Salvo i risultati in rows
	rows, err = db_Connection.Query(q, cod)
	//Se c'è un errore, ritorna un libro vuoto e l'errore
	if err != nil {
		return "", err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	//Creo una stringa in cui salvare il risultato della query
	var assegnatario string
	//Rows.Next() scorre tutte le righe trovate dalla query returnando true. Quando le finisce returna false
	for rows.Next() {
		//Tramite rows.Scan() salvo i vari risultati nel libro creato in precedenza. In caso di errore ritorno un libro vuoto e l'errore
		if err := rows.Scan(&assegnatario); err != nil {
			return "", err
		}
	}
	return assegnatario, err

}

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

//Funzione per trovare tutti i prestiti di un utente
func GetCurrentPrestito(codice uint32) (Prestito, error) {
	//Verifico se il server è ancora disponibile
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
		return Prestito{}, err
	}

	q := `SELECT COUNT(*) FROM Prestito WHERE Libro = ? AND Data_restituzione IS NULL`
	//Applico la query al database. Salvo i risultati in rows
	rows, err := db_Connection.Query(q, codice)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return Prestito{}, err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var count uint32
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return Prestito{}, err
		}
	}

	if count == 0 {
		return Prestito{}, &NoCurrentPrestitoError{codice}
	}

	//Salvo la query che eseguirà l'sql in una variabile stringa
	q = `SELECT Codice,Libro,Utente,Data_prenotazione,Durata FROM Prestito WHERE Libro = ? AND Data_restituzione IS NULL`
	//Applico la query al database. Salvo i risultati in rows
	rows, err = db_Connection.Query(q, codice)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return Prestito{}, err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var pres Prestito
	//Rows.Next() scorre tutte le righe trovate dalla query returnando true. Quando le finisce returna false
	for rows.Next() {
		var data_pre int64
		//Tramite rows.Scan() salvo i vari risultati nelle variabili create in precedenza. In caso di errore ritorno null e l'errore
		if err := rows.Scan(&pres.Codice, &pres.Libro, &pres.Utente, &data_pre, &pres.Durata); err != nil {
			return Prestito{}, err
		}
		//Salvo data_pre in pres convertendola in timestamp unix
		pres.Data_prenotazione = time.Unix(data_pre, 0)
	}
	//Se c'è un errore, ritorna null e l'errore
	if err := rows.Err(); err != nil {
		return Prestito{}, err
	}

	//Returno i prestiti trovati e null (null sarebbe l'errore che non è avvenuto)
	return pres, nil
}

//Funzione per la ricerca dei libri
func RicercaLibri(nome string, autore, genere []uint32, page int16) ([]Libro, error) {
	//Verifico se il server è ancora disponibile
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
		return nil, err
	}

	//Divido la stringa nome in vari tag e poi aggiungo i vari argomenti alla slice "args"
	tags := strings.Split(nome, " ")
	var args []interface{}
	for _, tag := range tags {
		args = append(args, "%"+tag+"%")
	}
	for _, a := range autore {
		args = append(args, a)
	}
	for _, g := range genere {
		args = append(args, g)
	}
	// Dividi per pagine solo se la pagina è un numero positivo
	if page >= 0 {
		args = append(args, uint16(page)*config.Config.Generale.LunghezzaPagina, uint16(page+1)*config.Config.Generale.LunghezzaPagina)
	}

	//Esamino tutti i casi possibili di richiesta, scegliendo la query giusta per ogni situazione possibile
	q := `SELECT Libro.Codice,Titolo,Autore.Nome,Autore.Cognome,Genere.Nome,Prenotato,Hashz FROM Libro,Autore,Genere WHERE Libro.Autore = Autore.Codice AND Libro.Genere = Genere.Codice` + strings.Repeat(` AND titolo LIKE ?`, len(tags))
	if len(autore) > 0 {
		q += ` AND autore IN (?` + strings.Repeat(`,?`, len(autore)-1) + `)`
	}
	if len(genere) > 0 {
		q += ` AND genere IN (?` + strings.Repeat(`,?`, len(genere)-1) + `)`
	}
	// Dividi per pagine solo se la pagina è un numero positivo
	if page >= 0 {
		q += ` LIMIT ?,?`
	}

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
		if err := rows.Scan(&fabrizio.Codice, &fabrizio.Titolo, &fabrizio.NomeAutore, &fabrizio.CognomeAutore, &fabrizio.Genere, &fabrizio.Prenotato, &fabrizio.Hashz); err != nil {
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
	for _, tag := range tags {
		if len(tag) > 0 {
			//I % servono per dire all'SQL di cercare la stringa in qualsiasi posizione
			args = append(args, "%"+tag+"%")
			args = append(args, "%"+tag+"%")
		}
	}

	q := `SELECT * FROM Autore WHERE 0 = 0` + strings.Repeat(` AND (nome LIKE ? OR cognome LIKE ?)`, len(args)/2)
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
	for _, tag := range tags {
		if len(tag) > 0 {
			//I % servono per dire all'SQL di cercare la stringa in qualsiasi posizione
			args = append(args, "%"+tag+"%")
		}
	}
	q := ``
	if len(args) > 0 {
		q = `SELECT * FROM Genere WHERE nome LIKE ?` + strings.Repeat(` OR nome LIKE ?`, len(args)-1)
	} else {
		q = `SELECT * FROM Genere`
	}
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

func AddPrestito(libro uint32, utente string, durata uint32) (uint32, error) {
	data_prenotazione := time.Now()
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

	q := `UPDATE Libro SET prenotato = 1 WHERE codice = ?`
	rows, err := db_Connection.Query(q, libro)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

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
func SetRestituzione(libro uint32) error {
	prestito_struct, err := GetCurrentPrestito(libro)
	if err != nil {
		return err
	}
	prestito := prestito_struct.Codice

	data_restituzione := time.Now()
	//Verifico se il server è ancora disponibile
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
		return err
	}
	q := `UPDATE Prestito
		  SET Data_restituzione = ?
		  WHERE codice = ?`
	rows, err := db_Connection.Query(q, data_restituzione.Unix(), prestito)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	q2 := `UPDATE Libro SET prenotato = 0 WHERE codice = ?`
	rows2, err := db_Connection.Query(q2, prestito)
	if err != nil {
		return err
	}
	defer rows2.Close()

	return nil
}

//Funzione per rimuovere un libro
func RemoveLibro(codice uint32) error {
	//Verifico se il server è ancora disponibile
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
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
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
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
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
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
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
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

//Funzione che ritorna il numero di libri prenotati
func LibriPrenotati() (uint32, error) {
	//Verifico se il server è ancora disponibile
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
		return 0, err
	}

	q := `SELECT COUNT(*) FROM Libro
		WHERE prenotato = 1`
	rows, err := db_Connection.Query(q)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return 0, err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var pren uint32
	for rows.Next() {
		//Tramite rows.Scan() salvo i vari risultati nella variabile creata in precedenza. In caso di errore ritorno null e l'errore
		if err := rows.Scan(&pren); err != nil {
			return 0, err
		}
	}

	//se c'è un errore, ritorna null e l'errore
	if err := rows.Err(); err != nil {
		return 0, err
	}

	return pren, nil
}

//Funzione che ritorna il numero di libri disponibili
func LibriDisponibili() (uint32, error) {
	//Verifico se il server è ancora disponibile
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
		return 0, err
	}

	q := `SELECT COUNT(*) FROM Libro
		WHERE prenotato = 0`
	rows, err := db_Connection.Query(q)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return 0, err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var disp uint32
	for rows.Next() {
		//Tramite rows.Scan() salvo i vari risultati nella variabile creata in precedenza. In caso di errore ritorno null e l'errore
		if err := rows.Scan(&disp); err != nil {
			return 0, err
		}
	}

	//se c'è un errore, ritorna null e l'errore
	if err := rows.Err(); err != nil {
		return 0, err
	}

	return disp, nil
}

//Funzione che conta gli autori totali
func CountAutori() (uint32, error) {
	//Verifico se il server è ancora disponibile
	//Se c'è un errore, ritorna null e l'errore
	if err := db_Connection.Ping(); err != nil {
		return 0, err
	}

	q := `SELECT COUNT(*) FROM Autore`
	rows, err := db_Connection.Query(q)
	//Se c'è un errore, ritorna null e l'errore
	if err != nil {
		return 0, err
	}
	//Rows verrà chiuso una volta che tutte le funzioni normali saranno terminate oppure al prossimo return
	defer rows.Close()

	var tot uint32
	for rows.Next() {
		//Tramite rows.Scan() salvo i vari risultati nella variabile creata in precedenza. In caso di errore ritorno null e l'errore
		if err := rows.Scan(&tot); err != nil {
			return 0, err
		}
	}

	//se c'è un errore, ritorna null e l'errore
	if err := rows.Err(); err != nil {
		return 0, err
	}

	return tot, nil

}
