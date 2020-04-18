/*
 * tipi.go
 *
 * Pacchetto per interfacciarsi con il database SQL
 *
 * Copyright (c) 2020 Filippo Casarin <casarin.filippo17@gmail.com>
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
