package masterpasswordmanager

import (
	"math/rand"
	"strings"
	"time"
)

var (
	lowerCaseSet        = "abcdefghijklmnopqrstuvwxyz"
	upperCaseSet        = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharacterSet = "!@#$%&*(){}[]?"
	digitsSet           = "0123456789"
	allCharacterSet     = lowerCaseSet + upperCaseSet + specialCharacterSet + digitsSet
)

const minAnyCharNumbers = 2
const masterpasswordLength = 25

// Generates a Crytographically secure random alphanumeric
func GenerateAccountSecretKey() string {

	rand.Seed(time.Now().UnixNano())

	var password strings.Builder

	for i := 0; i < minAnyCharNumbers; i++ {
		random_lower := rand.Intn(len(lowerCaseSet))
		password.WriteString(string(lowerCaseSet[random_lower]))
		random_upper := rand.Intn(len(upperCaseSet))
		password.WriteString(string(upperCaseSet[random_upper]))
		random_special := rand.Intn(len(specialCharacterSet))
		password.WriteString(string(specialCharacterSet[random_special]))
		random_digit := rand.Intn(len(digitsSet))
		password.WriteString(string(digitsSet[random_digit]))
	}

	masterpasswordRemianingLength := masterpasswordLength - 4*minAnyCharNumbers
	for i := 0; i < masterpasswordRemianingLength; i++ {
		random_all := rand.Intn(len(allCharacterSet))
		password.WriteString(string(allCharacterSet[random_all]))
	}

	masterPassword := []rune(password.String())
	rand.Shuffle(len(masterPassword), func(i, j int) {
		masterPassword[i], masterPassword[j] = masterPassword[j], masterPassword[i]
	})

	return string(masterPassword)
}
