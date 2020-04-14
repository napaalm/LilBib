# LilBib (Lightweight Integrated Logistics for Book Indexing and Borrowing)

Sistema di gestione bibliotecaria.

# Come contribuire

## Setup
Per contribuire al codice è consigliato eseguire dalla directory principale del repository il seguente comando:
```bash
$ ./misc/setup.sh
```

## Compilazione
Per compilare ed eseguire è possibile usare `make`:
```bash
$ make build && make run
```
È possibile anche generare un binario di release:
```bash
$ make release
```

Questi comandi creano le directory `/build` e `/release`. Per eliminarle:
```bash
$ make clean
```

# Struttura del progetto

## Pagine

### /
Home del sito web: può contenere informazioni sul progetto ed eventuali statistiche.

### /libri
Elenco dei libri con ricerca server-side.
I risultati sono divisi in più pagine.
La ricerca opera su titolo, autore e genere.
Di default reindirizza a `/libri/0`

#### /libri/\<page\>
Restituisce la `page`-esima pagina.

### /libro
Reindirizza a `/libri`.

#### /libro/\<id\>
Dettagli sul libro `id`.
Se il libro è correntemente in prestito visualizza l'assegnatario corrente.

### /autori
Elenco degli autori. Di default reindirizza ad `/autori/a`.

#### /autori/\<iniziale\>
Elenco degli autori con iniziale `iniziale`.
In cima è presente una lista con i collegamenti a tutte le iniziali disponibili.
Reindirizza alla ricerca di `/libri` quando si clicca su un autore.

### /generi
Elenco dei generi.
Reindirizza alla ricerca di `/libri` quando si clicca su un genere.

### /login
Pagina di accesso all'area utente.
Utilizza il server LDAP per l'autenticazione, ritorna un token e reindirizza a `/utente`.

### /utente
Contiene informazioni sull'utente, come il nome utente e la storia dei prestiti.
È presente un link a `/prestito`.

### /prestito
Permette di scansionare o inserire il codice di uno o più libri per prenderli in prestito scegliendone la durata.
In caso non sia stato effettuato l'accesso verranno richieste le proprie credenziali, senza però restituire un token (caso d'uso: computer comune in biblioteca per prendere in prestito e restituire libri).

### /restituzione
Permette di restituire i libri in proprio possesso.
Funzionamento identico a `/prestito`.

## Tabelle SQL

### Libri
* codice
* titolo
* autore
* genere
* prestito
* hash

### Generi
* codice
* nome

### Autore
* id
* nome
* cognome

### Prestiti
* codice
* libro
* utente
* data_prenotazione
* data_restituzione

## Backend GO

### Packages
* sql
* ldap
*

### Tipi
* Libro

### Funzioni
```go
func GetLibro(id uint32) Libro
```


