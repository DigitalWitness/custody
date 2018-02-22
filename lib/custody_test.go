package custody

import (
	"testing"

	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"database/sql"
	"github.com/gtank/cryptopasta"
	_ "github.com/mattn/go-sqlite3"
	"github.gatech.edu/NIJ-Grant/custody/crypto"
	"github.gatech.edu/NIJ-Grant/custody/models"
	"log"
)

func FailTest(t *testing.T, err error, fmtstring string) {
	if err != nil {
		t.Fatalf(fmtstring, err)
	}
}

func CheckCount(t *testing.T, db models.XODB, query string, expected int) {
	var cnt int
	res, err := db.Query(query)
	if err != nil {
		t.Fatalf("Query error %s", query)
	}
	res.Scan(&cnt)
	if cnt > expected {
		t.Fatalf("Number of records isn't at least: %s, %d", query, expected)
	}
}

func setupdb(t *testing.T, path string) *DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		t.Fatal(err)
	}
	cdb := &DB{db}
	if err = cdb.Init(); err != nil {
		t.Fatal(err)
	}
	return cdb
}

func TestDB_NewUser(t *testing.T) {

	cdb := setupdb(t, "./testing.sqlite")
	pub := []byte("BEGIN ECSDA KEY")
	type args struct {
		name      string
		publickey []byte
	}
	tests := []struct {
		name    string
		db      *DB
		args    args
		want    models.Identity
		wantErr bool
	}{
		{"evan", cdb, args{"evan", pub}, models.Identity{ID: 1, Name: "evan", CreatedAt: XONow(), PublicKey: pub}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.NewUser(tt.args.name, tt.args.publickey)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			want := tt.want
			if got.Name != want.Name {

				t.Errorf("DB.NewUser() = %v, want %v", got, tt.want)
			}
			if !bytes.Equal(got.PublicKey, got.PublicKey) {
				t.Fatalf("Keys not stored right")
			}
		})
	}
	CheckCount(t, cdb, "select count(*) from identities", 1)
	CheckCount(t, cdb, "select count(*) from ledger", 1)
}

func TestSignValidate(t *testing.T) {
	cdb := setupdb(t, "testing.sqlite")
	key, err := cryptopasta.NewSigningKey()
	if err != nil {
		t.Fatalf("Could not generate cryptokey %s", err)
	}
	data := []byte("This is a test")
	sig, err := cdb.Sign(data, key)
	if err != nil {
		t.Fatalf("Could not sign message %s", err)
	}
	valid := cryptopasta.Verify(data, sig, &key.PublicKey)
	if !valid {
		t.Fatalf("Validatation failed: %s, %v, %+v", data, sig, key.PublicKey)
	}
	data[0] = 'S'
	valid = cryptopasta.Verify(data, sig, &key.PublicKey)
	if valid {
		t.Fatalf("Validatation false positive: %s, %v, %+v", data, sig, key.PublicKey)
	}
	data[0] = 'T'
	sig[0] = sig[0] + 1
	valid = cryptopasta.Verify(data, sig, &key.PublicKey)
	if valid {
		t.Fatalf("Validatation false positive: %s, %v, %+v", data, sig, key.PublicKey)
	}
}

func TestLedger(t *testing.T) {
	cdb := setupdb(t, "./testing.sqlite")
	key, err := cryptopasta.NewSigningKey()
	data := []byte("Test of operation")
	hash, err := cryptopasta.Sign(data, key)
	FailTest(t, err, "failed to generate key %s.")
	pubbytes, err := x509.MarshalPKIXPublicKey(key.Public())
	FailTest(t, err, "failed to encode key %s.")
	i, err := cdb.NewUser("keyeduser", pubbytes)
	FailTest(t, err, "failed to create user %s.")
	ledg, err := cdb.Operate(&i, string(data), hash)

	FailTest(t, err, "failed to create ledger item %s.")
	if ledg.Identity != i.ID {
		t.Fatal(err)
	}
	if ledg.Message != "Test of operation" {
		t.Fatal("Failed to insert message correctly")
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

	pub, err = crypto.ParseECDSAPublicKey(keybytes)
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

func TestValidate(t *testing.T) {
	var err error
	var ident *models.Identity
	var keybytes []byte
	var pub *ecdsa.PublicKey
	cdb := setupdb(t, "testing.sqlite")

	req := Request{Command: "create"}
	ident, err = models.IdentityByID(cdb, 1)
	if err != nil {
		t.Fatalf("Could not find identity: %v", err)
	}
	req.Identity = *ident
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
	req.Identity.PublicKey = keybytes
	err = req.Identity.Update(cdb)
	FailTest(t, err, "could not update identity with public key: %v")

	id2, err := models.IdentityByID(cdb, 1)
	FailTest(t, err, "could not load identity 1")
	pub, err = id2.Public()
	FailTest(t, err, "could not load identity 1")

	message := []byte("Can you validate me?")

	sig, err := cryptopasta.Sign(message, key)
	FailTest(t, err, "Could not sign Message: %s")
	t.Logf("Sig: %s", sig)
	if pub == nil {
		t.Fatalf("Nil Pubkey")
	}
	valid := cryptopasta.Verify(message, sig, pub)
	if !valid {
		t.Fatalf("Could not validate message with key: %s", pub)
	}
}

//Can we get the identity entries associated with a username?
func TestIndexes(t *testing.T) {
	cdb := setupdb(t, "testing.sqlite")
	ids, err := models.IdentitiesByName(cdb, "evan")
	FailTest(t, err, "failed IdentitiesByName %s")
	for _, id := range ids {
		log.Print(*id)
	}
}

//Can we get the ledger entries associated with an identity?
func TestLedgerIndexes(t *testing.T) {
	cdb := setupdb(t, "testing.sqlite")
	ids, err := models.IdentitiesByName(cdb, "evan")
	id := ids[len(ids)-1]

	ls, err := models.LedgersByName(cdb, "evan")
	FailTest(t, err, "failed IdentitiesByName %s")
	for _, l := range ls {
		log.Print(*l)
		if l.Identity != id.ID {
			t.Fatal(l, id)
		}
	}
}
