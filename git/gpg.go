package git

import "zgit/setting"

const (
	FirstCommitScene = iota
)

func GetGpnKeyId(repo string, sceneType int) string {
	signKey := setting.SignKey()
	if signKey == "" || signKey == "default" {
		key, _ := GetRepoSignKey(repo)
		return key
	}
	return signKey
}
