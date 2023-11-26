package setting

import "github.com/LeeZXin/zsf/property/static"

var (
	defaultBranch = static.GetString("repo.defaultBranch")
)

func init() {
	if defaultBranch == "" {
		defaultBranch = "master"
	}
}

func DefaultBranch() string {
	return defaultBranch
}
