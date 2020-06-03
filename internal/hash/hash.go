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
	codice := uint32(0)
	for i := uint8(0); i < 4; i++ {
		codice += uint32(pass[i]) << (i * 8)
	}

	libro, err := db.GetLibro(codice)
	if err != nil {
		return libro, err
	}

	pass_decoded, err := base64.StdEncoding.DecodeString(pass)
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

	for i := uint8(0); codice != 0; i++ {
		pass[i] = byte(codice & 0xFF)
		codice /= 256
	}

	if _, err := rand.Read(pass[4:]); err != nil {
		return "", err
	}

	pass_encoded := base64.StdEncoding.EncodeToString(pass)

	hash := sha256.Sum256(pass[:])
	db.SetHash(codice, base64.StdEncoding.EncodeToString(hash[:]))

	return pass_encoded, nil
}
