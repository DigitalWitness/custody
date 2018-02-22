package client

import (
	"crypto/ecdsa"
	"crypto/x509"
	"github.com/mitchellh/go-homedir"
	"github.gatech.edu/NIJ-Grant/custody/crypto"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

//KeyDir: finds the path to the directoy containing the keys relative to fpath
// if path is empty, we default to homedirectory/.custodyctl/
func KeyDir(path string) (fpath string, err error) {
	if len(path) == 0 {
		path, err = homedir.Dir()
		path = filepath.Join(path, ".custodyctl")
	}
	fpath = path
	return
}

//LoadPublicKey: parse the public key from the base directory at dir,
//returns an error if we fail to read the x509 formatted file, or
//fail to parse the cert itself.
func LoadPublicKey(dir string) (key *ecdsa.PublicKey, err error) {
	path := dir
	pubpath := filepath.Join(path, "id_ecdsa.pub")
	keybytes, err := ioutil.ReadFile(pubpath)
	if err != nil {
		return
	}
	key, err = crypto.ParseECDSAPublicKey(keybytes)
	return
}

//LoadPrivateKey: parse the public key from the base directory at dir,
//returns an error if we fail to read the x509 formatted file, or
//fail to parse the cert itself.
func LoadPrivateKey(dir string) (key *ecdsa.PrivateKey, err error) {
	path := dir
	pubpath := filepath.Join(path, "id_ecdsa")
	keybytes, err := ioutil.ReadFile(pubpath)
	if err != nil {
		return
	}
	key, err = x509.ParseECPrivateKey(keybytes)
	return
}

//StoreKeys: write the public and private key-pair in x509 format to a directory.
//for example StoreKeys(key, "/home/user/.custodyctl/id_ecdsa") will store
//the keys as "/home/user/.custodyctl/id_ecdsa" and "/home/user/.custodyctl/id_ecdsa.pub"
//This naming convention is inspired by ssh-keygen.
//if path is empty, then store in $HOME/.custodyctl, if path is "./" then store in current directory.
func StoreKeys(key *ecdsa.PrivateKey, path string) (err error) {
	path, err = KeyDir(path)
	if err != nil {
		return
	}
	if err = os.MkdirAll(path, 0700); err != nil {
		return
	}
	log.Printf("writing keys to directory %s", path)

	privpath := filepath.Join(path, "id_ecdsa")
	fp, err := os.OpenFile(privpath, os.O_CREATE|os.O_RDWR, 0600)
	defer fp.Close()
	if err != nil {
		return
	}
	privbytes, err := x509.MarshalECPrivateKey(key)
	nb, err := fp.Write(privbytes)
	log.Printf("wrote Private key at path=%s, nbytes=%d", privpath, nb)
	if err != nil {
		return
	}
	pubpath := filepath.Join(path, "id_ecdsa.pub")
	fp, err = os.OpenFile(pubpath, os.O_CREATE|os.O_RDWR, 0644)
	defer fp.Close()
	if err != nil {
		return
	}
	pubbytes, err := x509.MarshalPKIXPublicKey(key.Public())
	nb, err = fp.Write(pubbytes)
	log.Printf("wrote Public key at path=%s, nbytes=%d", pubpath, nb)
	if err != nil {
		return
	}
	return

}
