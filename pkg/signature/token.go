package signature

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/LeeZXin/zsf-utils/strutil"
	"strings"
)

type User struct {
	Account string
	Email   string
}

func GetToken(user User) string {
	randomStr := strutil.RandomStr(8)
	return randomStr + getToken(randomStr, user)
}

func VerifyToken(token string, user User) bool {
	if len(token) < 8 {
		return false
	}
	extra := token[:8]
	sig := token[8:]
	return sig == getToken(extra, user)
}

func getToken(extra string, user User) string {
	h := sha256.New()
	_, _ = h.Write([]byte(strings.Join(
		[]string{
			user.Account,
			user.Email,
			extra,
		}, ":",
	)))
	return hex.EncodeToString(h.Sum(nil))
}
