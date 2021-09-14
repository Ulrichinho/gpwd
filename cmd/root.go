/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

// https://levelup.gitconnected.com/exploring-go-packages-cobra-fce6c4e331d6

var (
	cfgFile        string
	length         int
	quantity       int
	noSpecialsChar bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gpwd",
	Version: "1.0.1",
	Short:   "generate password(s)",
	Long:    `Golang CLI app which generate random password(s) with API : https://www.motdepasse.xyz/api`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("--length:", length)
		fmt.Println("--quantity:", quantity)
		fmt.Println("--no-specials-char:", noSpecialsChar)

		getRandomPassword(length, quantity, noSpecialsChar)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gpwd.yaml)")
	rootCmd.PersistentFlags().IntVarP(&length, "length", "l", 12, "define the length of password")
	rootCmd.PersistentFlags().IntVarP(&quantity, "quantity", "q", 1, "define the number of password to generate")
	rootCmd.PersistentFlags().BoolVar(&noSpecialsChar, "no-specials-char", false, "define if you don't want special character")
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

const APIURL = "https://api.motdepasse.xyz/create/"

type Password struct {
	Password []string `json:"passwords"`
}

func getRandomPassword(l, q int, nsc bool) {
	url := ""
	if nsc {
		url = APIURL + "?include_digits&include_lowercase&include_uppercase&password_length=" + strconv.Itoa(l) + "&quantity=" + strconv.Itoa(q)
	} else {
		url = APIURL + "?include_digits&include_lowercase&include_uppercase&include_special_characters&password_length=" + strconv.Itoa(l) + "&quantity=" + strconv.Itoa(q)
	}
	resBytes := getPasswordData(url)
	password := Password{}

	if err := json.Unmarshal(resBytes, &password); err != nil {
		fmt.Printf("Could not unmarshal reponseBytes. %v", err)
	}

	for _, p := range password.Password {
		fmt.Println(p)
	}
}

func getPasswordData(baseAPI string) []byte {
	r, err := http.NewRequest(
		http.MethodGet,
		baseAPI,
		nil,
	)

	if err != nil {
		log.Printf("Couldn't request a password. %v", err)
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("User-Agent", "https://www.motdepasse.xyz/api")

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Printf("Couln't make a request. %v", err)
	}

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Couln't read response body. %v", err)
	}

	return resBytes
}
