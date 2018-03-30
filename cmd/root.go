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
	"encoding/json"
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile, dsn string
var username string
var serverAddress string

// Config: the cmd configuration struct
type Config struct {
	cfgFile, dsn, username, serverAddress string
	json                                  bool
}

var config Config

var joutput *json.Encoder

// Output: print out an arbitrary value encoding using json
func Output(obj interface{}) {
	joutput = json.NewEncoder(os.Stdout)
	if config.json {
		joutput.Encode(obj)
	}
}

// Fatal: if err != nil, log.Fatal with a message.
func Fatal(err error, fmtstring string) {
	if err != nil {
		log.Fatalf(fmtstring, err)
	}
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "custody",
	Short: "A brief description of your application",
	Long: `Custody manager command line application. This application serves both as the client and server.
You can host a custody server using the serve command which will accept API requests to conduct operations on the database.
The operations include create, sign, and list which enroll users, sign messages or list application state.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.custody.yaml)")
	RootCmd.PersistentFlags().StringVar(&dsn, "dsn", "", "connection string for example file://custody.sqlite")
	RootCmd.PersistentFlags().StringVar(&username, "username", "", "the username for your identity")
	RootCmd.PersistentFlags().Bool("json", false, "use json formatted output")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	viper.BindPFlags(RootCmd.PersistentFlags())
	viper.BindPFlags(RootCmd.Flags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".custody" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".custody")
	}

	viper.SetEnvPrefix("CUST")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	viper.RegisterAlias("user", "username")
	dsn = viper.GetString("dsn")
	username = viper.GetString("username")
	serverAddress = "localhost"
	config.json = viper.GetBool("json")
	if config.json {
		log.Printf("using JSON output\n")
	} else {
		log.Printf("not using JSON output\n")
	}
	//log.Printf("Settings:\n%v\n", viper.AllSettings())
}
