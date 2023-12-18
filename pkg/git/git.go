package git

import (
	"bytes"
	"context"
	"fmt"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"zgit/pkg/git/command"
	"zgit/setting"
)

const (
	ZeroCommitId = "0000000000000000000000000000000000000000"
)

const (
	EnvRepoID     = "ZGIT_REPO_ID"
	EnvRepoPath   = "ZGIT_REPO_PATH"
	EnvRepoIsWiki = "ZGIT_REPO_IS_WIKI"
	EnvPusherID   = "ZGIT_PUSHER_ID"
	EnvPRID       = "ZGIT_PR_ID"
	EnvIsInternal = "ZGIT_INTERNAL_PUSH"
	EnvAppUrl     = "ZGIT_APP_URL"
)

const (
	RequiredVersion = "2.0.0"
	BranchPrefix    = "refs/heads/"
	TagPrefix       = "refs/tags/"
	TimeLayout      = "Mon Jan _2 15:04:05 2006 -0700"
	PrettyLogFormat = "--pretty=format:%H"
)

const notRegularFileMode = os.ModeSymlink | os.ModeNamedPipe | os.ModeSocket | os.ModeDevice | os.ModeCharDevice | os.ModeIrregular

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

func SetLocalConfig(ctx context.Context, repoPath, k, v string) error {
	_, err := command.NewCommand("config", "--local", k, v).Run(ctx, command.WithDir(repoPath))
	if err != nil {
		return fmt.Errorf("failed to set local config %s, err: %w", k, err)
	}
	return nil
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

// IsReferenceExist returns true if given reference exists in the repository.
func IsReferenceExist(ctx context.Context, repoPath, name string) bool {
	_, err := command.NewCommand("show-ref", "--verify", "--", name).Run(ctx, command.WithDir(repoPath))
	return err == nil
}

// IsBranchExist returns true if given branch exists in the repository.
func IsBranchExist(ctx context.Context, repoPath, name string) bool {
	if !strings.HasPrefix(name, BranchPrefix) {
		name = BranchPrefix + name
	}
	return IsReferenceExist(ctx, repoPath, name)
}

func HashObject(ctx context.Context, repoPath string, reader io.Reader) (string, error) {
	result, err := command.NewCommand("hash-object", "-w", "--stdin").
		Run(ctx, command.WithDir(repoPath), command.WithStdin(reader))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.ReadAsString()), nil
}

func AddObjectToIndex(ctx context.Context, repoPath, mode, object, filename string) error {
	_, err := command.NewCommand("update-index", "--add", "--replace", "--cacheinfo", mode, object, filename).
		Run(ctx, command.WithDir(repoPath))
	return err
}

// RemoveFilesFromIndex removes given filenames from the index - it does not check whether they are present.
func RemoveFilesFromIndex(ctx context.Context, repoPath string, filenames ...string) error {
	buffer := new(bytes.Buffer)
	for _, file := range filenames {
		if file != "" {
			buffer.WriteString("0 0000000000000000000000000000000000000000\t")
			buffer.WriteString(file)
			buffer.WriteByte('\000')
		}
	}
	_, err := command.NewCommand("update-index", "--remove", "-z", "--index-info").
		Run(ctx, command.WithDir(repoPath), command.WithStdin(bytes.NewReader(buffer.Bytes())))
	return err
}

// WriteTree writes the current index as a tree to the object db and returns its hash
func WriteTree(ctx context.Context, repoPath string) (*Tree, error) {
	result, err := command.NewCommand("write-tree").Run(ctx, command.WithDir(repoPath))
	if err != nil {
		return nil, err
	}
	return NewTree(strings.TrimSpace(result.ReadAsString())), nil
}

type CommitTreeOpts struct {
	Parents []string
	Message string
}

// CommitTree creates a commit from a given tree id for the user with provided message
func CommitTree(ctx context.Context, repoPath string, tree *Tree, opts CommitTreeOpts) (string, error) {
	cmd := command.NewCommand("commit-tree", tree.Id, "--no-gpg-sign")
	for _, parent := range opts.Parents {
		cmd.AddArgs("-p", parent)
	}
	message := new(bytes.Buffer)
	message.WriteString(opts.Message)
	message.WriteString("\n")
	result, err := cmd.Run(ctx, command.WithDir(repoPath), command.WithStdin(message))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.ReadAsString()), nil
}

func GetRepoSize(path string) (int64, error) {
	var size int64
	err := filepath.WalkDir(path, func(_ string, info os.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) { // ignore the error because the file maybe deleted during traversing.
				return nil
			}
			return err
		}
		if info.IsDir() {
			return nil
		}
		f, err := info.Info()
		if err != nil {
			return err
		}
		if (f.Mode() & notRegularFileMode) == 0 {
			size += f.Size()
		}
		return err
	})
	return size, err
}

type RefName string

func (n RefName) IsBranch() bool {
	return strings.HasPrefix(string(n), BranchPrefix)
}

func (n RefName) IsTag() bool {
	return strings.HasPrefix(string(n), TagPrefix)
}
