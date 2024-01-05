package setting

import (
	"errors"
	"github.com/LeeZXin/zsf-utils/idutil"
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
	dataDir, homeDir, appPath, repoDir string

	tempDir, lfsDir, avatarDir string

	appUrl = strings.TrimSuffix(static.GetString("app.url"), "/")

	isWindows = runtime.GOOS == "windows"

	resourcesDir string

	lang = static.GetString("app.lang")

	standaloneCorpId = static.GetString("app.corpId")

	hookToken = idutil.RandomUuid()
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		logger.Logger.Panicf("zgit os.Getwd err: %v", err)
	}
	dataDir = filepath.Join(pwd, "data")
	homeDir = filepath.Join(dataDir, "home")
	repoDir = filepath.Join(dataDir, "repo")
	tempDir = filepath.Join(dataDir, "temp")
	lfsDir = filepath.Join(dataDir, "lfs")
	avatarDir = filepath.Join(dataDir, "avatar")
	err = os.MkdirAll(homeDir, os.ModePerm)
	if err != nil {
		logger.Logger.Panicf("zgit os.MkdirAll homeDir err: %v", err)
	}
	err = os.MkdirAll(repoDir, os.ModePerm)
	if err != nil {
		logger.Logger.Panicf("zgit os.MkdirAll repoDir err: %v", err)
	}
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		logger.Logger.Panicf("zgit os.MkdirAll tempDir err: %v", err)
	}
	err = os.MkdirAll(lfsDir, os.ModePerm)
	if err != nil {
		logger.Logger.Panicf("zgit os.MkdirAll lfsDir err: %v", err)
	}
	err = os.MkdirAll(avatarDir, os.ModePerm)
	if err != nil {
		logger.Logger.Panicf("zgit os.MkdirAll avatarDir err: %v", err)
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

func RepoDir() string {
	return repoDir
}

func ResourcesDir() string {
	return resourcesDir
}

func TempDir() string {
	return tempDir
}

func LfsDir() string {
	return lfsDir
}

func Lang() string {
	return lang
}

func StandaloneCorpId() string {
	return standaloneCorpId
}

func HookToken() string {
	return hookToken
}
