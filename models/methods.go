package models

import (
	"crypto/ecdsa"
	"github.gatech.edu/NIJ-Grant/custody/crypto"
)

//Public: return the public key from an identity by
// parsing the x509 cert associated with the identity.
func (i *Identity) Public() (*ecdsa.PublicKey, error){
	return crypto.ParseECDSAPublicKey(i.PublicKey)
}
