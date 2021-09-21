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
package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const APIURL = "https://api.motdepasse.xyz/create/?"

type Password struct {
	Password []string `json:"passwords"`
}

// get API url in function of flags
func getAPIUrl(l, q int, nsc bool) string {
	if nsc {
		return APIURL + "include_digits&include_lowercase&include_uppercase&password_length=" + strconv.Itoa(l) + "&quantity=" + strconv.Itoa(q)
	}
	return APIURL + "include_digits&include_lowercase&include_uppercase&include_special_characters&password_length=" + strconv.Itoa(l) + "&quantity=" + strconv.Itoa(q)
}

// recover data from API
func getPasswordDataFromAPI(baseAPI string) []byte {
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

// recover data and format result
func GetRandomPassword(len, qa int, nsc bool) Password {
	url := getAPIUrl(len, qa, nsc)
	resBytes := getPasswordDataFromAPI(url)
	password := Password{}

	if err := json.Unmarshal(resBytes, &password); err != nil {
		fmt.Printf("Could not unmarshal reponseBytes. %v", err)
	}

	return password
}
