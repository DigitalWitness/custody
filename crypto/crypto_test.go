package crypto

import (
	"crypto/ecdsa"
	"github.com/gtank/cryptopasta"
	"crypto/x509"
	"testing"
)

func FailTest(t *testing.T, err error, fmtstring string) {
	if err != nil {
		t.Fatalf(fmtstring, err)
	}
}

func TestX509(t *testing.T) {
	var err error
	var keybytes []byte
	var pub *ecdsa.PublicKey

	key, err := cryptopasta.NewSigningKey()
	if err != nil {
		t.Fatalf("Could not generate key: %v", err)
	}
	//t.Logf("PriKey: %+v", key)
	//t.Logf("PubKey: %+v", key.PublicKey)

	keybytes, err = x509.MarshalPKIXPublicKey(key.Public())

	if err != nil {
		t.Fatal(err)
	}

	pub, err = ParseECDSAPublicKey(keybytes)
	FailTest(t, err, "could not restore key from bytes %v")
	message := []byte("Can you validate me?")

	sig, err := cryptopasta.Sign(message, key)
	FailTest(t, err, "Could not sign Message: %s")
	//t.Logf("Sig: %s", sig)
	if pub == nil {
		t.Fatalf("Nil Pubkey")
	}
	valid := cryptopasta.Verify(message, sig, pub)
	if !valid {
		t.Fatalf("Could not validate message with key: %s", pub)
	}
}
