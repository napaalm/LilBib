# LilBib (Lightweight Integrated Logistics for Book Indexing and Borrowing)

Sistema di gestione bibliotecaria.

# Come contribuire

## Setup
### Prerequisiti
È necessario installare [golang](https://golang.org/) e `make` per compilare il codice, `docker` e `docker-compose` per eseguire il database di sviluppo, e `zip` per comprimere i file di rilascio.

**Nota per gli utenti Windows**: si consiglia di installare [git for windows](https://gitforwindows.org/) poiché fornisce un set di tool simile ai sistemi Unix, il che è fondamentale al fine di eseguire il Makefile ed i git hooks. Inoltre, per installare `make` e `zip` si consiglia un package manager come [Chocolatey](https://chocolatey.org/).

### Inizializzazione progetto

Per contribuire al codice è necessario eseguire dalla directory principale del repository il seguente comando:
```bash
$ ./scripts/setup.sh
```

## Compilazione
Per compilare ed eseguire è possibile usare `make`:
```bash
$ make run
```
È possibile anche generare un binario di release:
```bash
$ make release
```

Per eliminare i risultati della compilazione:
```bash
$ make clean
```

# Organizzazione delle cartelle
Basata su [questo](https://github.com/golang-standards/project-layout).

### `/cmd/lilbib`
Contiene il main.

### `/config`
Contiene il file di configurazione di default.

### `/database`
Contiene la struttura del database ed un file docker-compose per avviare il database di sviluppo.

### `/githooks`
Contiene degli hook utili allo sviluppatore.

### `/internal`
Contiene i packages interni da cui è composto il software.

### `/scripts`
Contiene script utili allo sviluppatore.

### `/web`
Contiene gli asset statici e i template HTML.

### `/sandbox`
Contiene un ambiente di prova per i risultati della compilazione. (generata da `make`)

### `/release`
Contiene i binari di rilascio. (generata da `make`)

# Struttura del progetto

## Pagine e endpoint

### `/`
Home del sito web: può contenere informazioni sul progetto ed eventuali statistiche.

### `/libri`
Elenco dei libri con ricerca server-side.
I risultati sono divisi in più pagine.
La ricerca opera su titolo, autore e genere.
Di default reindirizza a `/libri/0`

#### `/libri/<page>`
Restituisce la `page`-esima pagina.

### `/libro`
Reindirizza a `/libri`.

#### `/libro/<id>`
Dettagli sul libro `id`.
Se il libro è correntemente in prestito visualizza l'assegnatario corrente.

### `/autori`
Elenco degli autori. Di default reindirizza ad `/autori/a`.

#### `/autori/<iniziale>`
Elenco degli autori con iniziale `iniziale`.
In cima è presente una lista con i collegamenti a tutte le iniziali disponibili.
Reindirizza alla ricerca di `/libri` quando si clicca su un autore.

### `/generi`
Elenco dei generi.
Reindirizza alla ricerca di `/libri` quando si clicca su un genere.

### `/login`
Pagina di accesso all'area utente.
Utilizza il server LDAP per l'autenticazione, ritorna un token e reindirizza a `/utente`.

### `/logout`
Endpoint che rimuove il token di accesso e reindirizza a `/login`.

### `/utente`
Contiene informazioni sull'utente, come il nome utente e la storia dei prestiti.
È presente un link a `/prestito`.

### `/prestito`
Permette di scansionare o inserire il codice di uno o più libri per prenderli in prestito scegliendone la durata.
~~In caso non sia stato effettuato l'accesso verranno richieste le proprie credenziali, senza però restituire un token (caso d'uso: computer comune in biblioteca per prendere in prestito e restituire libri).~~
In caso non sia stato effettuato l'accesso reindirizza a `/login`.

### `/restituzione`
Permette di restituire i libri in proprio possesso.
Funzionamento identico a `/prestito`.

### `/admin/generaCodici`
Permette all'utente amministratore di generare i codici QR per i libri.

### `/admin/aggiungiLibro`
Permette all'utente amministratore di aggiungere generi, autori e libri al database.

### `/static`
Endpoint per servire contenuti statici da `web/static`.

### `/api/getLibro?qrcode=<base64-encoded code+password>`
Endpoint per ottenere informazioni su un libro in formato JSON a partire dal contenuto del suo QR code (necessaria previa autenticazione).

### `/api/prestito?qrcode=<base64-encoded code+password>&durata=<time in seconds>`
Endpoint per prendere in prestito un libro per una certa durata (necessaria previa autenticazione).

### `/api/restituzione?qrcode=<base64-encoded code+password>`
Endpoint per restituire un libro (necessaria previa autenticazione).

## Tabelle SQL

### Libro
* `Codice` (primary key)
* `Titolo`
* `Autore`
* `Genere`
* `Prenotato`
* `Hashz`

### Genere
* `Codice` (primary key)
* `Nome`

### Autore
* `Codice` (primary key)
* `Nome`
* `Cognome`

### Prestito
* `Codice` (primary key)
* `Libro` (foreign key)
* `Utente` (foreign key)
* `Data_prenotazione`
* `Durata`
* `Data_restituzione`

## Backend GO
In ogni package è presente un file README dove sono indicate le funzionalità da implementare, i tipi e le funzioni esportate.
Si consiglia di scrivere funzioni interne al package per evitare che queste diventino troppo lunghe.

### Packages
* `auth`
* `config`
* `db`
* `handlers`
* `hash`
* `qrcode`

### Template codice sorgente
Questo template va incluso all'inizio di ogni file.
```go
/*
 * nome-file.go
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

```
