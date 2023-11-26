package setting

import (
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"os/exec"
)

var (
	gitUser = static.GetString("git.user")

	gitExecutablePath string
)

func init() {
	if gitUser == "" {
		gitUser = "git"
	}
	absPath, err := exec.LookPath("git")
	if err != nil {
		logger.Logger.Panicf("could not LookPath err: %v", err)
	}
	gitExecutablePath = absPath
}

func GitUser() string {
	return gitUser
}

func GitExecutablePath() string {
	return gitExecutablePath
}
