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
	"github.com/spf13/cobra"
	"github.gatech.edu/NIJ-Grant/custody/lib"
	"log"
	"net"
	"net/rpc"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"os/signal"
	"syscall"
	"time"
)



// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start the custodyctl server",
	Long: `The server must be running in order to conduct operations on the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("serve called")
		db, err := custody.Dial(dsn)
		if err != nil {
			log.Fatal(err)
		}
		cdb := custody.DB{db}
		log.Println(cdb)
		c := custody.NewClerk()
		c.DB = cdb
		rpc.Register(c)
		rpc.HandleHTTP()
		l, e := net.Listen(c.Network, c.Address)
		if e != nil {
			log.Fatal("listen error:", e)
		}
		go http.Serve(l, nil)


		var gracefulStop = make(chan os.Signal)
		signal.Notify(gracefulStop, syscall.SIGTERM)
		signal.Notify(gracefulStop, syscall.SIGINT)
		func() {
			sig := <-gracefulStop
			log.Printf("caught sig: %+v", sig)
			log.Println("Shutting down server")
			time.Sleep(500*time.Microsecond)
			os.Exit(0)
		}()

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
