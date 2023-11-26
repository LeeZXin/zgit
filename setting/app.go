package setting

import (
	"errors"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/zsf"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	dataDir, homeDir, appPath, repoRootDir string

	appUrl = static.GetString("app.url")

	isWindows = runtime.GOOS == "windows"

	resourcesDir string
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		logger.Logger.Panicf("zgit os.Getwd err: %v", err)
	}
	dataDir = filepath.Join(pwd, "data")
	homeDir = filepath.Join(dataDir, "home")
	repoRootDir = filepath.Join(dataDir, "repo")
	err = os.MkdirAll(homeDir, os.ModePerm)
	if err != nil {
		logger.Logger.Panicf("zgit os.MkdirAll homeDir err: %v", err)
	}
	err = os.MkdirAll(repoRootDir, os.ModePerm)
	if err != nil {
		logger.Logger.Panicf("zgit os.MkdirAll repoRootDir err: %v", err)
	}
	path, err := getAppPath()
	if err != nil {
		logger.Logger.Panicf("zgit getAppPath err: %v", err)
	}
	appPath = path
	rd, err := filepath.Abs("resources")
	if err != nil {
		logger.Logger.Panicf("zgit filepath.Abs(\"resources\") err: %v", err)
	}
	resourcesDir = rd
}

func getAppPath() (string, error) {
	var path string
	var err error
	if IsWindows() && filepath.IsAbs(os.Args[0]) {
		path = filepath.Clean(os.Args[0])
	} else {
		path, err = exec.LookPath(os.Args[0])
	}
	if err != nil {
		if !errors.Is(err, exec.ErrDot) {
			return "", err
		}
		path, err = filepath.Abs(os.Args[0])
	}
	if err != nil {
		return "", err
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(path, "\\", "/"), err
}

func DataDir() string {
	return dataDir
}

func HomeDir() string {
	return homeDir
}

func IsWindows() bool {
	return isWindows
}

func AppPath() string {
	return appPath
}

func AppUrl() string {
	return appUrl
}

func IsDebugRunMode() bool {
	return zsf.GetRunMode() == "debug"
}

func RepoRootDir() string {
	return repoRootDir
}

func ResourcesDir() string {
	return resourcesDir
}
