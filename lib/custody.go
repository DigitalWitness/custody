package custody

import (
	"github.gatech.edu/NIJ-Grant/custody/models"
	"time"
	"github.com/xo/xoutil"
	"crypto/ecdsa"
	"github.com/gtank/cryptopasta"
)

type DB struct {
	models.XODB
}

type Request struct {
	Command string `json:"command"`
	models.Identity
	models.Ledger
	ecdsa.PublicKey `json:"public_key"`
}

func (db *DB) Init() error {
	query := `

create table if not exists identities (
  id integer not null primary key,
  name text not null,
  created_at timestamp not null,
  public_key blob not null -- an x509 cert as ascii
);

create table if not exists ledger (
  id integer not null primary key,
  created_at timestamp not null,
  identity integer not null,
  message text not null,
  hash blob not null,

  foreign key (identity) references identities(id)
);
`
	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}

func XONow() xoutil.SqTime {
	return xoutil.SqTime{Time: time.Now()}
}

func (db *DB) NewUser(name string, publickey []byte) (models.Identity, error){
	t:=XONow()
	ident := models.Identity{Name:name, PublicKey:publickey, CreatedAt: t}
	if err := ident.Insert(db); err != nil {
		return ident, err
	}
	return ident, nil
}

func (db *DB) Operate(identity models.Identity, message string) (models.Ledger, error) {
	ledg := models.Ledger{Identity: identity.ID, Message:message}
	if err := ledg.Insert(db); err != nil {
		return ledg, err
	}
	return ledg, nil
}

func (db *DB) Validate(identity models.Identity, data []byte, hash []byte) (bool, error) {
	pub, err := identity.Public()
	if err != nil {
		return false, err
	}
	valid := cryptopasta.Verify(data, hash, pub)
	return valid, nil
}

func (db *DB) Sign(data []byte, key *ecdsa.PrivateKey) ([]byte, error) {
	sig, err := cryptopasta.Sign(data, key)
	if err != nil {
		return nil, err
	}
	return sig, nil
}
