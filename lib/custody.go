package custody

import (
	"crypto/ecdsa"
	"database/sql"
	"fmt"
	"github.com/gtank/cryptopasta"
	"github.com/xo/xoutil"
	"github.gatech.edu/NIJ-Grant/custody/models"
	"time"
)

type CustodyError struct {
	Operation  string
	ID         *models.Identity
	Message    []byte
	Signature  []byte
	publickey  *ecdsa.PublicKey
	privatekey *ecdsa.PrivateKey
}

func (e CustodyError) Error() string {
	return fmt.Sprintf("CryptoError: op:%s, id:%v", e.Operation, e.ID)
}

type DB struct {
	models.XODB
}

//Dial: connect to the custody server and return a handle to the connection.
//dsn argument describes the connection parameters
//TODO: change this to use an API layer.
func Dial(dsn string) (*DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	conn := &DB{db}
	conn.Init()
	return conn, nil
}

//Request: a structure for dispatching the network requests RPC style.
type Request struct {
	Command string `json:"command"`
	models.Identity
	models.Ledger
	ecdsa.PublicKey `json:"public_key"`
}

//Init: initialize the database by executing create table statements
//theses create tables use IF NOT EXISTS so that it will be idempotent
//if you change the schema, you must update this function and drop all
//the tables in the database or otherwise migrate the schema
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

//XONow: wrap the current time in an xoutil.SqTime so that it can be entered into the DB.
func XONow() xoutil.SqTime {
	return xoutil.SqTime{Time: time.Now()}
}

func (db *DB) NewUser(name string, publickey []byte) (models.Identity, error) {
	t := XONow()
	ident := models.Identity{Name: name, PublicKey: publickey, CreatedAt: t}
	if err := ident.Insert(db); err != nil {
		return ident, err
	}
	return ident, nil
}

func (db *DB) Operate(identity *models.Identity, message string, hash []byte) (ledg models.Ledger, err error) {
	var key *ecdsa.PublicKey
	var valid bool
	key, err = identity.Public()
	if err != nil {
		return
	}
	data := []byte(message)
	//log.Printf("data  read from stdin: %v", data)
	valid = cryptopasta.Verify(data, hash, key)
	if !valid {
		err = CustodyError{Operation: "InvalidSignature", ID: identity, Message: data, Signature: hash}
		return
	}
	ledg = models.Ledger{Identity: identity.ID, CreatedAt: XONow(), Message: message, Hash: hash}
	if err = ledg.Insert(db); err != nil {
		return
	}
	return
}

func (db *DB) Validate(identity models.Identity, data []byte, hash []byte) (bool, error) {
	pub, err := identity.Public()
	if err != nil {
		return false, err
	}
	valid := cryptopasta.Verify(data, hash, pub)
	return valid, nil
}

//Sign: use the private key to sign a message and insert it into the db.
func (db *DB) Sign(data []byte, key *ecdsa.PrivateKey) ([]byte, error) {
	sig, err := cryptopasta.Sign(data, key)
	if err != nil {
		return nil, err
	}
	return sig, nil
}
