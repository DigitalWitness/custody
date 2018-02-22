package models

import (
	"crypto/ecdsa"
	"fmt"
	"github.gatech.edu/NIJ-Grant/custody/crypto"
)

//Public: return the public key from an identity by
// parsing the x509 cert associated with the identity.
func (i *Identity) Public() (*ecdsa.PublicKey, error) {
	return crypto.ParseECDSAPublicKey(i.PublicKey)
}

//LedgersByName: look up all the ledger entries associated with a name
func LedgersByName(db XODB, name string) (ls []*Ledger, err error) {
	ids, err := IdentitiesByName(db, name)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		return ls, fmt.Errorf("no identity with name %s", name)
	}
	id := ids[len(ids)-1]
	ls, err = LedgersByIdentity(db, id.ID)
	return
}
