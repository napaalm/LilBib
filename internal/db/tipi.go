package db

import (
	"time"
)

type Libro struct {
	codice    uint32
	titolo    string
	autore    uint32
	genere    uint32
	prenotato bool
	hash      string
}

type Genere struct {
	codice uint32
	nome   string
}

type Autore struct {
	codice  uint32
	nome    string
	cognome string
}

type Prestito struct {
	codice            uint32
	libro             uint32
	utente            string
	data_prenotazione time.Time
	durata            uint32
	data_restituzione time.Time
}
