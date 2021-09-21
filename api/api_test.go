package api

import (
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestGetAPIURL(t *testing.T) {
	url := getAPIUrl(12, 1, false)
	// check url
	assert.Equal(t, url, "https://api.motdepasse.xyz/create/?include_digits&include_lowercase&include_uppercase&include_special_characters&password_length=12&quantity=1")
}

func TestGetPasswordFromAPI(t *testing.T) {
	url := "https://api.motdepasse.xyz/create/?include_digits&include_lowercase&include_uppercase&password_length=1&quantity=1"
	_, status := getPasswordDataFromAPI(url)
	// check connection with API
	assert.Equal(t, status, 200)
}

func TestGetRandomPassword(t *testing.T) {
	assert := assert.New(t)
	password := GetRandomPassword(12, 1, false).Password[0]
	// check length
	assert.Equal(len(password), 12)
	// match if there are symbols
	hasSymbol := func(s string) bool {
		for _, letter := range s {
			if unicode.IsSymbol(letter) {
				return true
			}
		}
		return false
	}
	assert.Equal(hasSymbol(password), true)
}
