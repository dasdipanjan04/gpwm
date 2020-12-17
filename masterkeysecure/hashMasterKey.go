package masterkeysecure

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateMasterKeyHash(masterKey string) string {
	hasherEngine := sha256.New()
	hasherEngine.Write([]byte(masterKey))
	return hex.EncodeToString(hasherEngine.Sum(nil))
}
