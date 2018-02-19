package crypto

import (
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
)

//ParseECDSAPUblicKey calls x509.ParsePKIXPublicKey and ensures that the result is an ecdsa.PublicKey.
func ParseECDSAPublicKey (data []byte) (*ecdsa.PublicKey, error){
	pubbox, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, err
	}
	switch pub := pubbox.(type) {
	default:
		return nil, fmt.Errorf("could not parse 509 PKI as ECDSA public key")
	case *ecdsa.PublicKey:
		return pub, nil

	}
}