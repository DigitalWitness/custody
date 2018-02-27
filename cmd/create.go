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

	"crypto/ecdsa"
	"crypto/x509"
	"log"
	"net/rpc"

	"github.com/gtank/cryptopasta"
	"github.com/spf13/cobra"
	"github.gatech.edu/NIJ-Grant/custody/client"
	"github.gatech.edu/NIJ-Grant/custody/lib"
	"github.gatech.edu/NIJ-Grant/custody/models"
)

//SubmitUser: user the API connection to create a user based on the username and the public key.
func SubmitIdentity(user string, key *ecdsa.PublicKey) (i models.Identity, err error) {
	rpcclient, err := rpc.DialHTTP("tcp", serverAddress+":4911")
	Fatal(err, "dialing: %s")
	keybytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		log.Printf("Bailing before request: %v", err)
		return
	}

	var reply models.Identity
	req := &custody.RecordRequest{Name: user, PublicKey: keybytes}
	log.Printf("Requesting Creation: %v", req)

	// Synchronous call
	err = rpcclient.Call("Clerk.Create", req, &reply)
	Fatal(err, "RPC call failed: %s")
	fmt.Printf("new user created: %d, %s, %s", reply.ID, reply.Name, reply.CreatedAt)
	return
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user for the custody system",
	Long:  `Enrolls a new user in the system by generating their x509 cert.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		log.Printf("creating user: %s", username)
		key, err := cryptopasta.NewSigningKey()
		Fatal(err, "could not generate key: %s")
		err = client.StoreKeys(key, "")
		Fatal(err, "could not store keys: %s")
		SubmitIdentity(username, &key.PublicKey)
	},
}

func init() {
	RootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
