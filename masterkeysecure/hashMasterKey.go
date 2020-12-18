package masterkeysecure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"gpwm/internal/glogger"
	"io"
)

func GenerateMasterKeyHashSha256(password string) [sha256.Size]byte {
	return sha256.Sum256([]byte(password))
}

func EncryptMasterKeyAES(data []byte, password string) []byte {

	key := GenerateMasterKeyHashSha256(password)

	cipherblock, err := aes.NewCipher(key[:])
	if err != nil {
		glogger.Glog("masterkeysecure:EncryptMasterKeyAES ", err.Error())
	}

	gcm, err := cipher.NewGCM(cipherblock)
	if err != nil {
		glogger.Glog("masterkeysecure:EncryptMasterKeyAES ", err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		glogger.Glog("masterkeysecure:EncryptMasterKeyAES ", err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext
}

func DecryptAESMasterKey(data []byte, password string) string {
	key := GenerateMasterKeyHashSha256(password)

	cipherblock, err := aes.NewCipher(key[:])
	if err != nil {
		glogger.Glog("masterkeysecure:DecryptAESMasterKey ", err.Error())
	}

	gcm, err := cipher.NewGCM(cipherblock)
	if err != nil {
		glogger.Glog("masterkeysecure:DecryptAESMasterKey ", err.Error())
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		glogger.Glog("masterkeysecure:DecryptAESMasterKey ", err.Error())
	}

	return string(plaintext)
}
