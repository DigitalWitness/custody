package client

import (
	"github.com/gtank/cryptopasta"
	"testing"
)

func FailTest(t *testing.T, err error, fmtstring string) {
	if err != nil {
		t.Fatalf(fmtstring, err)
	}
}

func TestKeyDir(t *testing.T) {
	p, err := KeyDir("./")
	if p != ".custodyctl" {
		t.Fatal(p, err)
	}
	p, err = KeyDir("/")
	if p != "/.custodyctl" {
		t.Fatal(p, err)
	}
}

func TestStoreKeys(t *testing.T) {
	key, err := cryptopasta.NewSigningKey()
	FailTest(t, err, "could not generate key: %s")
	err = StoreKeys(key, "./")
	FailTest(t, err, "could not store keys: %s")
	pubkey, err := LoadPublicKey("./")
	if err != nil {
		t.Fatal(err)
	}
	privkey, err := LoadPrivateKey("./")
	if err != nil {
		t.Fatal(err)
	}
	if pubkey == nil {
		t.Fatal("Failed to load pubkey")
	}
	if privkey == nil {
		t.Fatal("Failed to load privkey")
	}
}
