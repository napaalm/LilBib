package hash

import (
	"crypto/sha256"
	"testing"
	"encoding/base64"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/db"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/config"
)

func TestAll(t *testing.T) {
	config.LoadConfig("../../config/config.toml")

	if err := db.InizializzaDB(); err != nil {
		t.Error(err)
		return
	}

	codice, err := db.AddLibro("libro1", 13, 14)
	if err != nil {
		t.Error(err)
		return
	}

	if codice == 0 {
		t.Error("codice: 0")
		return
	}

	hash, pass, err := Genera(codice)
	if err != nil {
		t.Error(err)
		return
	}

	if len(pass) != 20 {
		t.Errorf("len(pass): %d\n", len(pass))
		return
	}

	checksum := sha256.Sum256(pass)
	if str := base64.StdEncoding.EncodeToString(checksum[:]); str != hash {
		t.Errorf("%s != %s\n", str, hash)
		return
	}

	db.SetHash(codice, hash)

	libro, err := Verifica(pass)
	if err != nil {
		t.Error(err)
		return
	}

	if libro.Codice != codice {
		t.Errorf("Il codice non corrisponde: %d != %d", libro.Codice, codice)
		return
	}
}
