package setting

import "github.com/LeeZXin/zsf/property/static"

var (
	signKey             = static.GetString("sign.key")
	signWhenFirstCommit = static.GetBool("sign.firstCommit")
	signUsername        = static.GetString("sign.username")
	signEmail           = static.GetString("sign.email")
)

func init() {
	if signUsername == "" {
		signUsername = "zgit"
	}
	if signEmail == "" {
		signEmail = "zgit@fake.local"
	}
}

func SignUsername() string {
	return signUsername
}

func SignEmail() string {
	return signEmail
}

func SignWhenFirstCommit() bool {
	return signWhenFirstCommit
}

func SignKey() string {
	return signKey
}
