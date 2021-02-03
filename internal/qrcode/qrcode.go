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
	"git.napaalm.xyz/napaalm/LilBib/internal/db"
	"git.napaalm.xyz/napaalm/LilBib/internal/hash"
	qrcode "github.com/skip2/go-qrcode"
	"strings"
)

type QRCodeLibro struct {
	Codice uint32
	Titolo string
	QRCode []byte
}

// Genera il codice QR per il libro del codice indicato.
func CreateQRCode(id uint32) (QRCodeLibro, error) {

	// Trova il titolo
	libro, err := db.GetLibro(id)

	if err != nil {
		return QRCodeLibro{}, err
	}

	titolo := libro.Titolo

	// Genera la password
	password, err := hash.Genera(id)

	if err != nil {
		return QRCodeLibro{}, err
	}

	// Genera il QR code
	png, err := qrcode.Encode(password, qrcode.Medium, 200)

	if err != nil {
		return QRCodeLibro{}, err
	}

	return QRCodeLibro{id, titolo, png}, nil
}

// Genera una pagina HTML con i codici disposti in una griglia.
func GeneratePage(ids []uint32) (string, error) {

	// Genera i QR code
	var qrcodes []QRCodeLibro

	for _, id := range ids {
		qrcode, err := CreateQRCode(id)

		if err != nil {
			return "", err
		}

		qrcodes = append(qrcodes, qrcode)
	}

	// Crea la pagina
	var page strings.Builder

	// Aggiunge l'header
	page.WriteString("<!DOCTYPE html><html><head><title>Codici QR - Lilbib</title></head><body>")

	// Aggiunge una tabella con le immagini
	page.WriteString("<table width=\"100%\">")

	for i, qrcode := range qrcodes {
		// Aggiunge una riga
		if i%3 == 0 {
			page.WriteString("<tr>")
		}

		// Allinea al centro
		page.WriteString("<td align=\"center\" valign=\"center\">")

		// Codifica l'immagine in base64
		image_encoded := base64.StdEncoding.EncodeToString(qrcode.QRCode)

		// Aggiunge l'immagine
		page.WriteString("<img src=\"data:image/png;base64, " + image_encoded + "\" />")

		// Aggiunge il titolo e chiude il tag td
		page.WriteString("<br />" + qrcode.Titolo + "</td>")

		// Chiude la riga ogni 3 elementi e all'ultimo
		if i%3 == 2 || i == len(ids)-1 {
			page.WriteString("</tr>")
		}
	}

	// Chiude la tabella
	page.WriteString("</table>")

	// Chiude la pagina
	page.WriteString("</body></html>")

	return page.String(), nil
}
