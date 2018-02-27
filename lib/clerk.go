package custody

import (
	"fmt"
	"github.gatech.edu/NIJ-Grant/custody/models"
	"log"
)

//NetConfig: a struct to hold network configuration information
type NetConfig struct {
	Network string
	Address string
}

//NewNetConfig: create a new NetConfig with default configuration
func NewNetConfig() NetConfig {
	return NetConfig{"tcp", "0.0.0.0:4911"}
}

//Clerk: a struct to represent the global state of the custody application
type Clerk struct {
	DB DB
	NetConfig
}

//NewClerk: create a new Clerk with default configuration
func NewClerk() *Clerk {
	return &Clerk{NetConfig: NewNetConfig()}
}

//Create: ask the clerk to create a user
func (c *Clerk) Create(req *CreationRequest, reply *models.Identity) (err error) {
	i, err := c.DB.NewUser(req.Name, req.PublicKey)
	if err != nil {
		return
	}
	*reply = i
	return
}

//Validate: ask the clerk to validate a message
func (c *Clerk) Validate(req *RecordRequest, reply *models.Ledger) (err error) {
	log.Printf("accessing identities %v", req.Name)
	ids, err := models.IdentitiesByName(c.DB, req.Name)
	log.Printf("accessed identities %v", ids)
	if err != nil || len(ids) < 1 {
		err = fmt.Errorf("no identities found with username:%s, err:%s", req.Name, err)
		return
	}
	i := ids[len(ids)-1]
	ledg, err := c.DB.Operate(i, string(req.Data), req.Hash)
	if err != nil {
		return
	}
	*reply = ledg
	return
}

//List: ask the clerk to list the ledger entries associated with an identity
func (c *Clerk) List(req *ListRequest, reply *[]*models.Ledger) (err error) {
	ls, err := models.LedgersByName(c.DB, req.Name)
	*reply = ls
	return
}
