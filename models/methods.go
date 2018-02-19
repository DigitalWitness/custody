package models

import (
	"crypto/ecdsa"
	"github.gatech.edu/NIJ-Grant/custody/crypto"
)

func (i *Identity) Public() (*ecdsa.PublicKey, error){
	return crypto.ParseECDSAPublicKey(i.PublicKey)
}
