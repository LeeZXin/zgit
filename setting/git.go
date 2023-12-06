package setting

import (
	"github.com/LeeZXin/zsf/logger"
	"os/exec"
)

var (
	gitExecutablePath string
)

func init() {
	absPath, err := exec.LookPath("git")
	if err != nil {
		logger.Logger.Panicf("could not LookPath err: %v", err)
	}
	gitExecutablePath = absPath
}

func GitUser() string {
	return "git"
}

func GitExecutablePath() string {
	return gitExecutablePath
}
