package hash

import (
	"git.antonionapolitano.eu/napaalm/LilBib/internal/config"
	"git.antonionapolitano.eu/napaalm/LilBib/internal/db"
	"testing"
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

	pass, err := Genera(codice)
	if err != nil {
		t.Error(err)
		return
	}

	//if len(pass) != 20 {
	//	t.Errorf("len(pass): %d\n", len(pass))
	//	return
	//}

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
