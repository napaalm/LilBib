/*
 * hash.go
 *
 * Pacchetto per la generazione e verifica dei codici dei libri.
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

// Pacchetto per la generazione e verifica dei codici dei libri.
package hash

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/db"
)

type ErrHash struct{}

func (e ErrHash) Error() string {
	return "Book hash didn't match"
}

func Verifica(pass string) (db.Libro, error) {
	pass_decoded, err := base64.RawURLEncoding.DecodeString(pass)
	if err != nil {
		return db.Libro{}, err
	}

	codice := uint32(0)
	for i := uint8(0); i < 4; i++ {
		codice += uint32(pass_decoded[i]) << (i * 8)
	}

	libro, err := db.GetLibro(codice)
	if err != nil {
		return libro, err
	}

	hash := sha256.Sum256(pass_decoded)
	if base64.StdEncoding.EncodeToString(hash[:]) != libro.Hashz {
		return libro, ErrHash{}
	}

	return libro, nil
}

func Genera(codice uint32) (string, error) {
	pass := make([]byte, 20)

	temp := codice
	for i := uint8(0); temp != 0; i++ {
		pass[i] = byte(temp & 0xFF)
		temp /= 256
	}

	if _, err := rand.Read(pass[4:]); err != nil {
		return "", err
	}

	pass_encoded := base64.RawURLEncoding.EncodeToString(pass)

	hash := sha256.Sum256(pass[:])
	db.SetHash(codice, base64.StdEncoding.EncodeToString(hash[:]))

	return pass_encoded, nil
}
