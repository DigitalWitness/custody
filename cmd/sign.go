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

	"io/ioutil"
	"log"
	"net/rpc"
	"os"

	"github.com/gtank/cryptopasta"
	"github.gatech.edu/NIJ-Grant/custody/client"
	"github.gatech.edu/NIJ-Grant/custody/lib"
	"github.gatech.edu/NIJ-Grant/custody/models"
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "sign creates a ledger entry signed by the current user.",
	Long: `Signed entries can be used to record operations on files attributed to users.
You need the private key stored in ~/.custodyctl/id_ecdsa in order to create a valid signature.
The custody create command is used to generate key pairs and upload the public part to the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var reply models.Ledger
		fmt.Println("signing message from stdin")

		key, err := client.LoadPrivateKey("")
		Fatal(err, "could not load public key: %s")

		data, err := ioutil.ReadAll(os.Stdin)
		Fatal(err, "could not read input: %s")
		log.Printf("bytes read from stdin: %d", len(data))
		log.Printf("string read from stdin: %s", data)

		hash, err := cryptopasta.Sign(data, key)
		Fatal(err, "could not hash input: %s")
		//log.Printf("Successful hashing: %s", hash)

		client, err := rpc.DialHTTP("tcp", serverAddress+":4911")
		Fatal(err, "dialing: %s")

		req := custody.RecordRequest{Name: username, Data: data, Hash: hash}
		err = client.Call("Clerk.Validate", &req, &reply)
		Fatal(err, "could not add message to ledger %s")
		log.Printf("Ledger Entry: %+v", reply)
	},
}

func init() {
	RootCmd.AddCommand(signCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// signCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// signCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
