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

	"github.gatech.edu/NIJ-Grant/custody/client"
	"github.com/gtank/cryptopasta"
	"io/ioutil"
	"os"
	"log"
	"github.gatech.edu/NIJ-Grant/custody/lib"
	"github.gatech.edu/NIJ-Grant/custody/models"
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sign called")
		var err error
		db, err = custody.Dial(dsn)
		Fatal(err, "could not connect to API: %s")
		log.Printf("Database DSN=%s, DB=%+v", dsn, db)
		keydir, err := client.KeyDir("")
		Fatal(err, "could not find key dir")
		log.Printf("Keydir: %s", keydir)
		key, err := client.LoadPrivateKey(keydir)
		Fatal(err, "could not load public key: %s")
		data, err := ioutil.ReadAll(os.Stdin)
		Fatal(err, "could not read input: %s")
		log.Printf("bytes read from stdin: %d", len(data))
		log.Printf("string read from stdin: %v", data)
		hash, err := cryptopasta.Sign(data, key)
		Fatal(err,"could not hash input: %s")
		//log.Printf("Successful hashing: %s", hash)

		myid = 3
		i, err := models.IdentityByID(db, myid)
		ledg, err := db.Operate(i, string(data), hash)
		Fatal(err, "could not add message to ledger %s")
		log.Printf("Ledger Entry: %v", ledg)
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
