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
	"log"
	"github.com/gtank/cryptopasta"
	"crypto/x509"
)

func Fatal(err error, fmtstring string) {
	if err != nil {
		log.Fatalf(fmtstring, err)
	}
}
// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user for the custody system",
	Long: `Enrolls a new user in the system by generating their x509 cert.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
		fmt.Println(args)
		if len(args) < 1 {
			log.Fatal("Not enough arguments")
		}
		username := args[0]
		log.Printf("user: %s", username)
		key, err := cryptopasta.NewSigningKey()
		Fatal(err, "could not generate key: %s")
		keybytes, err := x509.MarshalPKIXPublicKey(key.Public())
		Fatal(err, "could not marshal public key: %s")
		log.Printf("Public Key: %s", keybytes)
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
