package hash

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/db"
)

type ErrHash struct{}

func (e ErrHash) Error() string {
	return "Book hash didn't match"
}

func Verifica(pass []byte) (db.Libro, error) {
	codice := uint32(0)
	for i := uint8(0); i < 4; i++ {
		codice = uint32(pass[i]) << (i * 8)
	}

	libro, err := db.GetLibro(codice)

	if err != nil {
		return libro, err
	}

	hash := sha256.Sum256(pass)

	if !bytes.Equal(hash[:], libro.hash) {
		return libro, ErrHash{}
	}

	return libro, nil
}

func Genera(codice uint32) ([]byte, []byte, error) { // hash password
	pass := make([]byte, 20)

	for i := uint8(0); codice != 0; i++ {
		pass[i] = byte(codice & 0xFF)
		codice /= 256
	}

	if _, err := rand.Read(pass[4:]); err != nil {
		return nil, nil, err
	}

	hash := sha256.Sum256(pass[:])

	return hash[:], pass, nil
}
