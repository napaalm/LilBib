/*
 * qrcode.go
 *
 * Funzioni per la generazione dei QR code destinati ai libri.
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

// Funzioni per la generazione dei QR code destinati ai libri.
package qrcode

import (
	"encoding/base64"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/db"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/hash"
	qrcode "github.com/skip2/go-qrcode"
)

type QRCodeLibro struct {
	Codice uint32
	Titolo string
	QRCode []byte
}

// Genera il codice QR per il libro del codice indicato.
func CreateQRCode(id uint32) (QRCodeLibro, error) {
	// Trova il titolo
	titolo, err := db.GetLibro(id).Titolo

	if err != nil {
		return QRCodeLibro{}, nil
	}

	// Genera la password
	password, err := hash.Genera(id)

	if err != nil {
		return QRCodeLibro{}, nil
	}

	// Genera il QR code
	png, err := qrcode.Encode(password, qrcode.Medium, 256)

	if err != nil {
		return QRCodeLibro{}, nil
	}

	return QRCodeLibro{id, titolo, png}
}
