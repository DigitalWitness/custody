// Copyright © 2018 James Fairbanks <james.fairbanks@gatech.edu>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.gatech.edu/NIJ-Grant/custody/lib"
	"log"
	"net"
	"net/rpc"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
	"github.gatech.edu/NIJ-Grant/custody/models"
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
	DB custody.DB
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
		err = fmt.Errorf("no identities found with username:%s, err:%s", username, err)
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


// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")
		db, err := custody.Dial(dsn)
		if err != nil {
			log.Fatal(err)
		}
		cdb := custody.DB{db}
		fmt.Println(cdb)
		c := NewClerk()
		c.DB = cdb
		rpc.Register(c)
		rpc.HandleHTTP()
		l, e := net.Listen(c.Network, c.Address)
		if e != nil {
			log.Fatal("listen error:", e)
		}
		http.Serve(l, nil)
	},
}


func init() {
	RootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
