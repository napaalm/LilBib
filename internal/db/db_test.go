/*
 * db_test.go
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
	"fmt"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/config"
	"testing"
)

func TestAll(t *testing.T) {
	config.LoadConfig("../../config/config.toml")

	if err := InizializzaDB(); err != nil {
		t.Error(err)
		return
	}

	defer ChiudiDB()

	genere, err := AddGenere("test_genere")
	if err != nil {
		t.Error(err)
		return
	}

	autore, err := AddAutore("test_nome", "test_cognome")
	if err != nil {
		t.Error(err)
		return
	}

	libro, err := AddLibro("test_libro", autore, genere)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = RicercaLibri("test libro", nil, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}

	if err := SetHash(libro, "assr"); err != nil {
		t.Error(err)
		return
	}

	if err := RemoveLibro(libro); err != nil {
		t.Error(err)
		return
	}

	if err := RemoveGenere(genere); err != nil {
		t.Error(err)
		return
	}

	if err := RemoveAutore(autore); err != nil {
		t.Error(err)
		return
	}

	if _, err := AddPrestito(1, "test_user", 10000); err != nil {
		t.Error(err)
		return
	}

	if err := SetRestituzione(1); err != nil {
		t.Error(err)
		return
	}

	//suqi, err := GetCurrentPrestito(2)
	//fmt.Println(GetCurrentPrestito(2))
	//err = SetRestituzione(2)
	//fmt.Println(GetCurrentPrestito(2))

	// id, err := db.AddPrestito(8, "marcoilbeffardo", time.Now(), 100)
	// auths, err := db.RicercaAutori("")
	// libro, err := db.GetLibro(1)
	// leib, err := db.RicercaLibri("inquisizione stalin chievo verona", nil, nil, 1, 3)
	// auts, err := db.GetAutori(97)
	// gen, err := db.GetGeneri()
	// prests, err := db.GetPrestiti("buffone02")
}
