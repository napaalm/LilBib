package hash

import (
	"fmt"
	"crypto/sha256"
	"testing"
	"encoding/base64"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/db"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/config"
)

func TestDatabase(t *testing.T) {
	config.LoadConfig("../../config/config.toml")

	_, err := db.GetLibro(1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestVerifica(t *testing.T) {
	pass, err := base64.StdEncoding.DecodeString("AQAAAPx/Y5uKUBfNWH8Cijwt55w=")
	if err != nil {
		t.Error(err)
	}

	if _, err = Verifica(pass); err != nil {
		t.Error(err)
	}
}

func TestGenera(t *testing.T) {
	hash, pass, err := Genera(0x1)
	if err != nil {
		t.Error(err)
	}

	// if len(hash) != 32 {
	// 	t.Errorf("len(hash): %d\n", len(hash))
	// }

	if len(pass) != 20 {
		t.Errorf("len(pass): %d\n", len(pass))
	}

	var xxx [100]int64
	for i := range xxx {
		xxx[i] = 0
	}

	checksum := sha256.Sum256(pass)
	if base64.StdEncoding.EncodeToString(checksum[:]) != hash {
		t.Errorf("%x != %x\n", checksum, hash)
	}

	fmt.Println(hash)
	fmt.Println(base64.StdEncoding.EncodeToString(pass))
}
