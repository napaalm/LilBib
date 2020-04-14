# LilBib

Sistema di prenotazione libri.

# Struttura del progetto
## Tabelle SQL
### Utenti?
* id
* Username
* Email?

### Libri
* id
* Titolo
* Autore
* Genere
* Prestito
* Hash

### Generi
* id
* Nome

### Autore
* id
* Nome
* Cognome

### Prestiti
* id
* Libro
* Utente
* Data prenotazione
* Data restituzione

## Backend GO
### Tipi
* Libro

### Funzioni
```go
func GetLibro(id uint32) Libro
```

## Pagine
### /
Più o meno niente

### /libri
Elenco dei libri con ricerca

### /generi
Elenco dei generi
Reidizza alla ricerca `/libri` quando si preme su un genere

### /autori/<iniziale>
Elenco degli autori con iniziale `iniziale`  
Reindirizza alla ricerca `/libri` quando si preme su un autore

### /libro/<id>
Dettagli sul libro `id`

* /login
* /prestiti
* /prenota
* /restituzione
