package masterkeysecure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"github.com/dasdipanjan04/gpwm/internal/glogger"
)

func GenerateMasterKeyHashSha256(password string) [sha256.Size]byte {
	return sha256.Sum256([]byte(password))
}

func EncryptMasterKeyAES(data []byte, password string) ([]byte, error) {

	key := GenerateMasterKeyHashSha256(password)

	cipherblock, err := aes.NewCipher(key[:])
	if err != nil {
		glogger.Glog("masterkeysecure:EncryptMasterKeyAES ", err.Error())
		return nil, err
	}

	gcm, err := cipher.NewGCM(cipherblock)
	if err != nil {
		glogger.Glog("masterkeysecure:EncryptMasterKeyAES ", err.Error())
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		glogger.Glog("masterkeysecure:EncryptMasterKeyAES ", err.Error())
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

func DecryptAESMasterKey(data []byte, password string) (string, error) {
	key := GenerateMasterKeyHashSha256(password)

	cipherblock, err := aes.NewCipher(key[:])
	if err != nil {
		glogger.Glog("masterkeysecure:DecryptAESMasterKey ", err.Error())
		return "", err
	}

	gcm, err := cipher.NewGCM(cipherblock)
	if err != nil {
		glogger.Glog("masterkeysecure:DecryptAESMasterKey ", err.Error())
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		glogger.Glog("masterkeysecure:DecryptAESMasterKey ", err.Error())
		return "", err
	}

	return string(plaintext), err
}
