package gpwmcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"github.com/dasdipanjan04/gpwm/helper/glogger"
)

// GenerateKEKHashSha256 generates a Key Encryption Key(KEK) using user provider secret password and some non-secret such as user account email address.
func GenerateKEKHashSha256(password string, email string) [sha256.Size]byte {
	return sha256.Sum256([]byte(password + email))
}

// EncryptKEKAES encrypts Actual Key By the KEK
func EncryptKEKAES(data []byte, password string, email string) ([]byte, error) {

	key := GenerateKEKHashSha256(password, email)

	cipherblock, err := aes.NewCipher(key[:])
	if err != nil {
		glogger.Glog("gpwmcrypto:EncryptMasterKEKAES ", err.Error())
		return nil, err
	}

	gcm, err := cipher.NewGCM(cipherblock)
	if err != nil {
		glogger.Glog("gpwmcrypto:EncryptMasterKEKAES ", err.Error())
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		glogger.Glog("gpwmcrypto:EncryptMasterKEKAES ", err.Error())
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

// Decrypt Actual key using KEK
func DecryptAESKEK(data []byte, password string, email string) (string, error) {
	key := GenerateKEKHashSha256(password, email)

	cipherblock, err := aes.NewCipher(key[:])
	if err != nil {
		glogger.Glog("gpwmcrypto:DecryptAESMasterKEK:Cipherblock ", err.Error())
		return "", err
	}

	gcm, err := cipher.NewGCM(cipherblock)
	if err != nil {
		glogger.Glog("gpwmcrypto:DecryptAESMasterKEK:GCM ", err.Error())
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		glogger.Glog("gpwmcrypto:DecryptAESMasterKEK:Plaintext ", err.Error())
		return "", err
	}

	return string(plaintext), err
}
