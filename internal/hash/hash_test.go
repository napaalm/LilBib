/*
 * hash_test.go
 *
 * File di test per il package hash.
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
 * <http://Www.gnu.org/licenses/>.
 */

package hash

import (
	"fmt"
	"git.napaalm.xyz/napaalm/LilBib/internal/config"
	"git.napaalm.xyz/napaalm/LilBib/internal/db"
	"testing"
)

func TestAll(t *testing.T) {
	config.LoadConfig("../../config/config.toml")

	if err := db.InizializzaDB(); err != nil {
		t.Error(err)
		return
	}

	//codice, err := db.AddLibro("libro1", 13, 14)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}

	//if codice == 0 {
	//	t.Error("codice: 0")
	//	return
	//}

	pass, err := Genera(1)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(pass)

	if len(pass) != 27 {
		t.Errorf("len(pass): %d\n", len(pass))
		return
	}

	libro, err := Verifica(pass)
	if err != nil {
		t.Error(err)
		return
	}

	if libro.Codice != 1 {
		t.Errorf("Il codice non corrisponde: %d != %d", libro.Codice, 1)
		return
	}
}
