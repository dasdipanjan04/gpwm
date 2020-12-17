package masterkeysecure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"log"
)

func GenerateMasterKeyHashSha256(masterKey string) [sha256.Size]byte {
	return sha256.Sum256([]byte(masterKey))
}

func EncryptMasterKeyAES(data []byte, masterKey string) []byte {

	key := GenerateMasterKeyHashSha256(masterKey)

	cipherblock, err := aes.NewCipher(key[:])
	if err != nil {
		log.Fatalln(err)
	}

	gcm, err := cipher.NewGCM(cipherblock)
	if err != nil {
		log.Fatalln(err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalln(err)
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext
}

func DecryptAESMasterKey(data []byte, masterkey string) string {
	key := GenerateMasterKeyHashSha256(masterkey)

	cipherblock, err := aes.NewCipher(key[:])
	if err != nil {
		log.Fatalln(err)
	}

	gcm, err := cipher.NewGCM(cipherblock)
	if err != nil {
		log.Fatalln(err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatalln(err)
	}

	return string(plaintext)
}
