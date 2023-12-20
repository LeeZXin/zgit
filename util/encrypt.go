package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func EncryptUserPassword(pwd string) string {
	h := sha256.New()
	h.Write([]byte(pwd))
	return hex.EncodeToString(h.Sum(nil))
}
