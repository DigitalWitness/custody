package custody

import (
	"fmt"
	"log"

	"github.gatech.edu/NIJ-Grant/custody/models"
)

// NetConfig: a struct to hold network configuration information
type NetConfig struct {
	Network string
	Address string
}

// NewNetConfig: create a new NetConfig with default configuration
func NewNetConfig() NetConfig {
	return NetConfig{"tcp", "0.0.0.0:4911"}
}

// Clerk: a struct to represent the global state of the custody application.
// The clerk is used to register functions for RPC.
// each method of the Clerk is accessible through the server using an RPC client.
type Clerk struct {
	DB DB
	NetConfig
}

// NewClerk: create a new Clerk with default configuration
func NewClerk() *Clerk {
	return &Clerk{NetConfig: NewNetConfig()}
}

// Create: ask the clerk to create a user
func (c *Clerk) Create(req *RecordRequest, reply *models.Identity) (err error) {
	if req.PublicKey == nil {
		err = fmt.Errorf("you must provide an x509 ECDSA public key with a user creation request")
	}
	i, err := c.DB.NewUser(req.Name, req.PublicKey)
	if err != nil {
		return
	}
	*reply = i
	return
}

// Validate: ask the clerk to validate a message
func (c *Clerk) Validate(req *RecordRequest, reply *models.Ledger) (err error) {
	var ids []*models.Identity
	var ledg models.Ledger
	log.Printf("clerk is accessing identities of user: %v", req.Name)
	ids, err = models.IdentitiesByName(c.DB, req.Name)
	if err != nil || len(ids) < 1 {
		err = fmt.Errorf("no identities found with username:%s, err:%s", req.Name, err)
		return
	}
	i := ids[len(ids)-1]
	ledg, err = c.DB.Operate(i, string(req.Data), req.Hash)
	if err != nil {
		return
	}
	*reply = ledg
	return
}

// List: ask the clerk to list the ledger entries associated with an identity
func (c *Clerk) List(req *RecordRequest, reply *[]*models.Ledger) (err error) {
	ls, err := models.LedgersByName(c.DB, req.Name)
	*reply = ls
	return
}
