# LilBib (Lightweight Integrated Logistics for Book Indexing and Borrowing)

Sistema di gestione bibliotecaria.

# Come contribuire

## Setup
Per contribuire al codice è consigliato eseguire dalla directory principale del repository il seguente comando:
```bash
$ ./scripts/setup.sh
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

Per eliminare i risultati della compilazione:
```bash
$ make clean
```

# Organizzazione delle cartelle
Basata su [questo](https://github.com/golang-standards/project-layout).

### `/cmd/lilbib`
Contiene il main.

### `/internal`
Contiene i packages interni da cui è composto il software.

### `/vendor`
Contiene le dipendenze del progetto.

### `/web`
Contiene gli asset statici e i template HTML.

### `/scripts`
Contiene script utili allo sviluppatore.

### `/release`
Contiene i binari di rilascio.

### `/githooks`
Contiene degli hook utili allo sviluppatore.

# Struttura del progetto

## Pagine

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

### `/utente`
Contiene informazioni sull'utente, come il nome utente e la storia dei prestiti.
È presente un link a `/prestito`.

### `/prestito`
Permette di scansionare o inserire il codice di uno o più libri per prenderli in prestito scegliendone la durata.
In caso non sia stato effettuato l'accesso verranno richieste le proprie credenziali, senza però restituire un token (caso d'uso: computer comune in biblioteca per prendere in prestito e restituire libri).

### `/restituzione`
Permette di restituire i libri in proprio possesso.
Funzionamento identico a `/prestito`.

### `/static`
Endpoint per servire contenuti statici da `web/static`.

## Tabelle SQL

### Libro
* `codice`
* `titolo`
* `autore`
* `genere`
* `prenotato`
* `hash`

### Genere
* `codice`
* `nome`

### Autore
* `codice`
* `nome`
* `cognome`

### Prestito
* `codice`
* `libro`
* `utente`
* `data_prenotazione`
* `durata`
* `data_restituzione`

## Backend GO
In ogni package è presente un file README dove sono indicate le funzionalità da implementare, i tipi e le funzioni esportate.
Si consiglia di scrivere funzioni interne al package per evitare funzioni troppo lunghe.

### Packages
* `auth`
* `config`
* `db`
* `handlers`
* `hash`

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
