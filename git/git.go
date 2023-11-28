package git

import (
	"fmt"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"zgit/git/command"
	"zgit/setting"
)

const (
	EnvRepoName     = "GITEA_REPO_NAME"
	EnvRepoUsername = "GITEA_REPO_USER_NAME"
	EnvRepoID       = "GITEA_REPO_ID"
	EnvRepoIsWiki   = "GITEA_REPO_IS_WIKI"
	EnvPusherName   = "GITEA_PUSHER_NAME"
	EnvPusherEmail  = "GITEA_PUSHER_EMAIL"
	EnvPusherID     = "GITEA_PUSHER_ID"
	EnvKeyID        = "GITEA_KEY_ID"
	EnvDeployKeyID  = "GITEA_DEPLOY_KEY_ID"
	EnvPRID         = "GITEA_PR_ID"
	EnvIsInternal   = "GITEA_INTERNAL_PUSH"
	EnvAppURL       = "GITEA_ROOT_URL"
	EnvActionPerm   = "GITEA_ACTION_PERM"
)

const (
	RequiredVersion = "2.0.0"
	BranchPrefix    = "refs/heads/"
	TagPrefix       = "refs/tags/"
	TimeLayout      = "Mon Jan _2 15:04:05 2006 -0700"
	PrettyLogFormat = "--pretty=format:%H"
	RemotePrefix    = "refs/remotes/"
	PullPrefix      = "refs/pull/"
)

var (
	supportProcReceive = CheckGitVersionAtLeast("2.29") == nil
)

func InitGit() {
	if CheckGitVersionAtLeast(RequiredVersion) != nil {
		logger.Logger.Panic("install git version is not supported, upgrade it before start zgit")
	}
	if _, ok := os.LookupEnv("GNUPGHOME"); !ok {
		_ = os.Setenv("GNUPGHOME", filepath.Join(setting.HomeDir(), ".gnupg"))
	}
	if CheckGitVersionAtLeast("2.18") == nil {
		command.AddGlobalCmdArgs("-c", "protocol.version=2")
	}
	if CheckGitVersionAtLeast("2.9") == nil {
		command.AddGlobalCmdArgs("-c", "credential.helper=")
	}
	if setting.LfsEnabled() {
		if CheckGitVersionAtLeast("2.1.2") != nil {
			logger.Logger.Panic("LFS server support requires Git >= 2.1.2")
		}
		command.AddGlobalCmdArgs("-c", "filter.lfs.required=", "-c", "filter.lfs.smudge=", "-c", "filter.lfs.clean=")
	}
	options := map[string]string{
		"diff.algorithm":  "histogram",
		"gc.reflogExpire": "90",

		"core.logAllRefUpdates": "true",
		"core.quotePath":        "false",
	}
	if CheckGitVersionAtLeast("2.10") == nil {
		options["receive.advertisePushOptions"] = "true"
	}
	if CheckGitVersionAtLeast("2.18") == nil {
		options["core.commitGraph"] = "true"
		options["gc.writeCommitGraph"] = "true"
		options["fetch.writeCommitGraph"] = "true"
	}
	if static.Exists("git.reflog.core.logAllRefUpdates") {
		options["core.logAllRefUpdates"] = strconv.FormatBool(static.GetBool("git.reflog.core.logAllRefUpdates"))
	}
	if static.GetInt("git.reflog.gc.reflogExpire") > 0 {
		options["gc.reflogExpire"] = static.GetString("git.reflog.gc.reflogExpire")
	}
	for k, v := range options {
		mustSetGlobalConfig(k, v)
	}
	mustSetGlobalConfigIfAbsent("user.name", setting.SignUsername())
	mustSetGlobalConfigIfAbsent("user.email", setting.SignEmail())
	if supportProcReceive {
		mustAddGlobalConfigIfAbsent("receive.procReceiveRefs", "refs/for")
	} else {
		mustUnsetAllGlobalConfig("receive.procReceiveRefs", "refs/for")
	}
	mustAddGlobalConfigIfAbsent("safe.directory", "*")
	if setting.IsWindows() {
		mustSetGlobalConfig("core.longpaths", "true")
		mustUnsetAllGlobalConfig("core.protectNTFS", "false")
	}
	if CheckGitVersionAtLeast("2.22") == nil {
		mustSetGlobalConfig("uploadpack.allowfilter", "true")
		mustSetGlobalConfig("uploadpack.allowAnySHA1InWant", "true")
	}
}

func mustSetGlobalConfig(k, v string) {
	if err := SetGlobalConfig(k, v); err != nil {
		logger.Logger.Panic(err)
	}
}

func mustSetGlobalConfigIfAbsent(k, v string) {
	if err := SetGlobalConfigIfAbsent(k, v); err != nil {
		logger.Logger.Panic(err)
	}
}

func mustAddGlobalConfigIfAbsent(k, v string) {
	if err := AddGlobalConfigIfAbsent(k, v); err != nil {
		logger.Logger.Panic(err)
	}
}

func mustUnsetAllGlobalConfig(k, v string) {
	if err := UnsetAllGlobalConfig(k, v); err != nil {
		logger.Logger.Panic(err)
	}
}

func SetGlobalConfigIfAbsent(k, v string) error {
	return setGlobalConfigCheckOverwrite(k, v, false)
}

func SetGlobalConfig(k, v string) error {
	return setGlobalConfigCheckOverwrite(k, v, true)
}

func setGlobalConfigCheckOverwrite(k, v string, overwrite bool) error {
	result, err := command.NewCommand("config", "--global", "--get").AddArgs(k).Run(nil)
	// fatal error
	if err != nil && !command.IsExitCode(err, 1) {
		return fmt.Errorf("failed to get git config %s, err: %w", k, err)
	}
	// 如果配置存在但不覆盖
	if err == nil && !overwrite {
		return nil
	}
	var currValue string
	// 配置存在
	if err == nil {
		currValue = strings.TrimSpace(result.ReadAsString())
	}
	if currValue == v {
		return nil
	}
	_, err = command.NewCommand("config", "--global").AddArgs(k, v).Run(nil)
	if err != nil {
		return fmt.Errorf("failed to set git global config %s, err: %w", k, err)
	}
	return nil
}

func AddGlobalConfigIfAbsent(k, v string) error {
	_, err := command.NewCommand("config", "--global", "--get").AddArgs(k, regexp.QuoteMeta(v)).Run(nil)
	if err == nil {
		return nil
	}
	if command.IsExitCode(err, 1) {
		_, err = command.NewCommand("config", "--global", "--add").AddArgs(k, v).Run(nil)
		if err != nil {
			return fmt.Errorf("failed to add git global config %s, err: %w", k, err)
		}
		return nil
	}
	return fmt.Errorf("failed to get git config %s, err: %w", k, err)
}

func UnsetAllGlobalConfig(k, v string) error {
	_, err := command.NewCommand("config", "--global", "--get").AddArgs(k).Run(nil)
	if err == nil {
		_, err = command.NewCommand("config", "--global", "--unset-all").AddArgs(k, regexp.QuoteMeta(v)).Run(nil)
		if err != nil {
			return fmt.Errorf("failed to unset git global config %s, err: %w", k, err)
		}
		return nil
	}
	if command.IsExitCode(err, 1) {
		return nil
	}
	return fmt.Errorf("failed to get git config %s, err: %w", k, err)
}

func getRefCompareSeparator(directCompare bool) string {
	if directCompare {
		return ".."
	}
	return "..."
}
