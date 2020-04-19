package hash

import (
	"bytes"
	"crypto/sha256"
	"testing"
)

func TestGenera(t *testing.T) {
	hash, pass, err := Genera(0x5428)

	if err != nil {
		t.Error(err)
	}

	if len(hash) != 32 {
		t.Errorf("len(hash): %d\n", len(hash))
	}

	if len(pass) != 20 {
		t.Errorf("len(pass): %d\n", len(pass))
	}

	var xxx [100]int64
	for i := range xxx {
		xxx[i] = 0
	}

	checksum := sha256.Sum256(pass[:])
	if !bytes.Equal(checksum[:], hash) {
		t.Errorf("%x != %x\n", checksum, hash)
	}
}
