/*
Copyright © 2021 GROLHIER Ulrich <grolhier.u@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"gpwd/api"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

// flag variables
var (
	cfgFile        string
	length         int
	quantity       int
	noSpecialsChar bool
	export         bool
	statistic      bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gpwd",
	Version: "1.0.1",
	Short:   "Golang CLI app which generate random password(s) with API : https://www.motdepasse.xyz/api",
	Long:    `This CLI app allows to generate passwords or keys with a limit of 512 characters in length and a maximum of 30 passwords at the same time. With the possibility of making exports and having an eye on statistics`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()

		switch {
		case length > 512:
			color.Cyan("[INFO] Cannot create password with more of 512 character for reasons of limit of API")
			os.Exit(1)
		case length < 12 && length >= 8:
			color.Yellow("[WARNING] it's not recommended to generate password(s) with a length less than 12 chars")
		case length < 8:
			color.Red("[ALERT] it's not secure!!Length min is of 8")
			os.Exit(1)
		}

		if quantity > 30 {
			color.Cyan("[INFO] Cannot create more of 30 passwords for reasons of limit of API")
			os.Exit(1)
		}

		passwords := api.GetRandomPassword(length, quantity, noSpecialsChar)

		switch {
		case export && len(args) == 0:
			WritePasswordInFile(args, &passwords)

			api.MapPassword(quantity, &passwords)

			color.Green("[SUCCESS] %v password(s) export in 'password.txt'\n", quantity)
		case export && len(args) == 1:
			WritePasswordInFile(args, &passwords)

			api.MapPassword(quantity, &passwords)

			color.Green("[SUCCESS] %v password(s) export in '%v.txt'\n", quantity, args[0])
		default:
			api.MapPassword(quantity, &passwords)
		}

		end := time.Now()

		if statistic {
			elapsed := end.Sub(start)
			fmt.Printf("Finished in %v\n", elapsed)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gpwd.yaml)")
	rootCmd.PersistentFlags().IntVarP(&length, "length", "l", 12, "define the length of password")
	rootCmd.PersistentFlags().IntVarP(&quantity, "quantity", "q", 1, "define the number of password to generate")
	rootCmd.PersistentFlags().BoolVar(&noSpecialsChar, "no-specials-char", false, "define if you don't want special character")
	rootCmd.PersistentFlags().BoolVarP(&export, "export", "e", false, "define if you want export passwords")
	rootCmd.PersistentFlags().BoolVarP(&statistic, "statistic", "s", false, "log the stats (speed, ...)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gpwd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gpwd")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func CheckErr(msg interface{}) {
	if msg != nil {
		fmt.Fprintln(os.Stderr, "Error:", msg)
		os.Exit(1) 
	}
}

// write password in a file
func WritePasswordInFile(args []string, passwords *api.Password) {
	var f *os.File
	var err error
	if len(args) == 0 {
		f, err = os.OpenFile("password.txt", os.O_CREATE|os.O_WRONLY, 0600)
		CheckErr(err)
	} else {
		f, err = os.OpenFile(args[0]+".txt", os.O_CREATE|os.O_WRONLY, 0600)
		CheckErr(err)
	}
	defer f.Close()

	for _, p := range passwords.Password {
		_, err = f.WriteString(p + "\n")
		CheckErr(err)
	}
}
