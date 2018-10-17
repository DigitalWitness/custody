package custody

import (
	"testing"

	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"database/sql"
	"log"

	"io/ioutil"

	"github.com/gtank/cryptopasta"
	_ "github.com/mattn/go-sqlite3"
	"github.gatech.edu/NIJ-Grant/custody/client"
	"github.gatech.edu/NIJ-Grant/custody/models"
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

func TestUploadExistingKey(t *testing.T) {

	//checkfounduser: helper function to test that the user we made exists
	checkfounduser := func(t *testing.T, cdb *DB, name string) {
		ids, err := models.IdentitiesByName(cdb, name)
		if err != nil {
			t.Fatalf("failed to find user from premade key %s", err.Error())
		}
		for i, id := range ids {
			if id.Name != name {
				t.Fatalf("found a bogus user %d, %v", i, id)
			}
		}
	}
	//createcheckuser: helper function to insert a user and then check that it exists
	createcheckuser := func(t *testing.T, cdb *DB, name string, pubkey []byte) {
		cdb.NewUser("premade_user", pubkey)
		checkfounduser(t, cdb, name)
	}

	// make a tempdir to store the keys in.
	tmpdir, err := ioutil.TempDir("", "custodyctl")
	t.Logf("Tempdir for keys is %s", tmpdir)
	cdb := setupdb(t, "testing.sqlite")
	var pubkey []byte

	// make a new key
	key, err := cryptopasta.NewSigningKey()
	if err != nil {
		t.Fatalf("failed to create key %s", err.Error())
	}

	// marshall the key into bytes "by hand"
	pubkey, err = x509.MarshalPKIXPublicKey(key.Public())
	if err != nil {
		t.Fatalf("failed to marshall the key %s", err.Error())
	}

	// assert that we can make a user having made the key externally
	createcheckuser(t, cdb, "premade_user", pubkey)

	// this simulates the creation of the keys by the android app
	err = client.StoreKeys(key, tmpdir)
	if err != nil {
		t.Fatalf("failed to store keys to filesystem %s", err)
	}

	// key gets loaded, you could read this file by hand
	// see the definition of client.LoadPublicKey for how to read a ECDSA key
	// into a key object.
	loadedkey, err := client.LoadPublicKey(tmpdir)
	if err != nil {
		t.Fatalf("failed to read back key from filesystem %s", err)
	}
	pubkey, err = x509.MarshalPKIXPublicKey(loadedkey)
	if err != nil {
		t.Fatalf("failed to marshall the key %s", err.Error())
	}
	createcheckuser(t, cdb, "premade_user_file", pubkey)

	// you can also insert premade users with a Clerk / with RPC
	ck := NewClerk()
	ck.DB = *cdb
	req := RecordRequest{Name: "premade_clerk", PublicKey: pubkey}
	reply := new(models.Identity)
	err = ck.Create(&req, reply)
	if err != nil {
		t.Fatalf("clerk failed to add premade user %s", err.Error())
	}
	checkfounduser(t, &ck.DB, "premade_clerk")
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

func TestValidate(t *testing.T) {
	var err error
	var ident *models.Identity
	var keybytes []byte
	var pub *ecdsa.PublicKey
	cdb := setupdb(t, "testing.sqlite")

	req := Request{Operation: Create}
	ident, err = models.IdentityByID(cdb, 1)
	if err != nil {
		t.Fatalf("Could not find identity: %v", err)
	}
	req.Identity = *ident
	key, err := cryptopasta.NewSigningKey()
	if err != nil {
		t.Fatalf("Could not generate key: %v", err)
	}
	// t.Logf("PriKey: %+v", key)
	// t.Logf("PubKey: %+v", key.PublicKey)

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

// Can we get the identity entries associated with a username?
func TestIndexes(t *testing.T) {
	cdb := setupdb(t, "testing.sqlite")
	ids, err := models.IdentitiesByName(cdb, "evan")
	FailTest(t, err, "failed IdentitiesByName %s")
	for _, id := range ids {
		log.Print(*id)
	}
}

// Can we get the ledger entries associated with an identity?
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
