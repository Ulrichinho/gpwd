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

	"golang.design/x/clipboard"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

// https://levelup.gitconnected.com/exploring-go-packages-cobra-fce6c4e331d6

var (
	cfgFile        string
	length         int
	quantity       int
	noSpecialsChar bool
	export         bool
)

var rootCmd = &cobra.Command{
	Use:     "gpwd",
	Version: "1.0.2",
	Short:   "generate password(s)",
	Long:    `Golang CLI app which generate random password(s) with API : https://www.motdepasse.xyz/api`,
	Run: func(cmd *cobra.Command, args []string) {
		// test length flag
		switch {
		case length < 12 && length >= 8:
			color.Yellow("[WARNING] it's not recommended to generate password(s) with a length less than 12 chars")
		case length < 8:
			color.Red("[ALERT] it's not secure!!Length min is of 8")
			os.Exit(1)
		}

		if quantity > 32 {
			color.Cyan("[INFO] Cannot create more of 32 passwords for reasons of performance")
			os.Exit(1)
		}

		getRandomPassword(length, quantity, noSpecialsChar, export)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gpwd.yaml)")
	rootCmd.PersistentFlags().IntVarP(&length, "length", "l", 12, "define the length of password")
	rootCmd.PersistentFlags().IntVarP(&quantity, "quantity", "q", 1, "define the number of password to generate")
	rootCmd.PersistentFlags().BoolVar(&noSpecialsChar, "no-specials-char", false, "define if you don't want special character")
	rootCmd.PersistentFlags().BoolVarP(&export, "export", "e", false, "define if you want export passwords")
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

func getUrl(l, q int, nsc bool) string {
	if nsc {
		return APIURL + "?include_digits&include_lowercase&include_uppercase&password_length=" + strconv.Itoa(l) + "&quantity=" + strconv.Itoa(q)
	}
	return APIURL + "?include_digits&include_lowercase&include_uppercase&include_special_characters&password_length=" + strconv.Itoa(l) + "&quantity=" + strconv.Itoa(q)
}

func mapPassword(q int, pwd *Password) {
	if q == 1 {
		for _, p := range pwd.Password {
			fmt.Println(p)
			clipboard.Write(clipboard.FmtText, []byte(p))
			fmt.Println("Copied")
		}
	} else {
		for _, p := range pwd.Password {
			fmt.Println(p)
		}
	}
}

func exportPassword(q int, e bool, pwd *Password) {
	if e {
		f, err := os.OpenFile("password.txt", os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		for _, p := range pwd.Password {
			_, err = f.WriteString(p + "\n")
			if err != nil {
				panic(err)
			}
		}

		mapPassword(q, pwd)

		color.Green("[SUCCESS] %v password(s) export in 'password.txt'\n", q)
	} else {
		mapPassword(q, pwd)
	}
}

// recover data and format result
func getRandomPassword(l, q int, nsc bool, e bool) {
	url := getUrl(l, q, nsc)
	resBytes := getPasswordData(url)
	password := Password{}

	if err := json.Unmarshal(resBytes, &password); err != nil {
		fmt.Printf("Could not unmarshal reponseBytes. %v", err)
	}

	// export ?
	exportPassword(q, e, &password)
}

// recover data from API
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
