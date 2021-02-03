/*
 * qrcode_test.go
 *
 * File di test per il package qrcode.
 *
 * Copyright (c) 2020 Antonio Napolitano <nap@antonionapolitano.eu>
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

package qrcode

import (
	"fmt"
	"git.napaalm.xyz/napaalm/LilBib/internal/config"
	"git.napaalm.xyz/napaalm/LilBib/internal/db"
	"testing"
)

func TestAll(t *testing.T) {
	// Carica e visualizza la configurazione di test
	config.LoadConfig("../../config/config.toml")
	fmt.Println(config.Config)

	// Carica il database
	db.InizializzaDB()

	// Libri di esempio
	var ids []uint32 = []uint32{1, 2, 4}

	// Genera la pagina dei QR code e stampala
	page, err := GeneratePage(ids)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(page)
	}
}
