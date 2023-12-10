package git

import "zgit/setting"

type CommitScene int

const (
	FirstCommitScene CommitScene = iota
)

func GetGpnKeyId(repo string, sceneType CommitScene) string {
	switch sceneType {
	case FirstCommitScene:
		if !setting.SignWhenFirstCommit() {
			return ""
		}
	default:
		break
	}
	signKey := setting.SignKey()
	if signKey == "" || signKey == "default" {
		key, _ := GetRepoSignKey(repo)
		return key
	}
	return signKey
}
