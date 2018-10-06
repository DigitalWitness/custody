package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"github.gatech.edu/NIJ-Grant/custody/lib"
)

// serveCmd represents the serve command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "start the custody web server",
	Long:  `The server must be running in order to conduct operations on the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("starting web server")
		server, _ := custody.InitializeHTTPServer()
		server.ListenAndServe()
	},
}

func init() {
	RootCmd.AddCommand(httpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}