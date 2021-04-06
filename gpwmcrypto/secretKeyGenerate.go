package gpwmcrypto

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

//GenerateAccountSecretKey Generates a Crytographically secure random alphanumeric
func GenerateAccountSecretKey() string {

	rand.Seed(time.Now().UnixNano())

	var password strings.Builder

	for i := 0; i < minAnyCharNumbers; i++ {
		randomLower := rand.Intn(len(lowerCaseSet))
		password.WriteString(string(lowerCaseSet[randomLower]))
		randomUpper := rand.Intn(len(upperCaseSet))
		password.WriteString(string(upperCaseSet[randomUpper]))
		randomSpecial := rand.Intn(len(specialCharacterSet))
		password.WriteString(string(specialCharacterSet[randomSpecial]))
		randomDigit := rand.Intn(len(digitsSet))
		password.WriteString(string(digitsSet[randomDigit]))
	}

	masterpasswordRemianingLength := masterpasswordLength - 4*minAnyCharNumbers
	for i := 0; i < masterpasswordRemianingLength; i++ {
		randomAll := rand.Intn(len(allCharacterSet))
		password.WriteString(string(allCharacterSet[randomAll]))
	}

	masterPassword := []rune(password.String())
	rand.Shuffle(len(masterPassword), func(i, j int) {
		masterPassword[i], masterPassword[j] = masterPassword[j], masterPassword[i]
	})

	return string(masterPassword)
}
