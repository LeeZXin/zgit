package git

import (
	"context"
	"errors"
	"fmt"
	"github.com/LeeZXin/zsf/logger"
	"os"
	"path/filepath"
	"strings"
	"time"
	"zgit/git/command"
	"zgit/setting"
	"zgit/util"
)

var (
	gitIgnoreResourcesPath = filepath.Join(setting.ResourcesDir(), "gitignore")
)

type Repository struct {
	Id    string `json:"id"`
	Owner User   `json:"owner"`
	Name  string `json:"name"`
	Path  string `json:"path"`
}

type InitRepoOpts struct {
	CreateReadme  bool
	GitIgnoreName string
	DefaultBranch string
}

func initEmptyRepository(ctx context.Context, repoPath string) error {
	if repoPath == "" {
		return errors.New("repoPath is empty")
	}
	logger.Logger.Infof("init repo: %s", repoPath)
	isExist, err := util.IsExist(repoPath)
	if err != nil {
		return err
	}
	if isExist {
		return fmt.Errorf("%s is exist", repoPath)
	}
	err = os.MkdirAll(repoPath, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = command.NewCommand("init", command.BareFlag).Run(ctx, command.WithDir(repoPath))
	return err
}

func InitRepository(ctx context.Context, repo Repository, opts InitRepoOpts) error {
	if err := initEmptyRepository(ctx, repo.Path); err != nil {
		return err
	}
	if opts.CreateReadme || opts.GitIgnoreName != "" {
		tmpDir, err := os.MkdirTemp(os.TempDir(), "zgit-"+repo.Name)
		if err != nil {
			return fmt.Errorf("failed to create temp dir for repository %s: %w", repo.Path, err)
		}
		defer func() {
			if err := util.RemoveAll(tmpDir); err != nil {
				logger.Logger.Errorf("Unable to remove temporary directory: %s: Error: %v", tmpDir, err)
			}
		}()
		return initTemporaryRepository(ctx, repo, tmpDir, opts)
	}
	branch := opts.DefaultBranch
	if branch == "" {
		branch = setting.DefaultBranch()
	}
	SetDefaultBranch(ctx, repo.Path, branch)
	return nil
}

func initTemporaryRepository(ctx context.Context, repo Repository, tmpDir string, opts InitRepoOpts) error {
	if err := CloneRepository(ctx, repo.Path, tmpDir, nil); err != nil {
		return fmt.Errorf("failed to clone original repository %s: %w", repo.Name, err)
	}
	if opts.CreateReadme {
		util.WriteFile(filepath.Join(tmpDir, "README.md"), []byte(fmt.Sprintf("# %s  ", repo.Name)))
	}
	if opts.GitIgnoreName != "" {
		content, err := os.ReadFile(filepath.Join(gitIgnoreResourcesPath, opts.GitIgnoreName))
		if err == nil {
			util.WriteFile(filepath.Join(tmpDir, ".gitIgnore"), content)
		}
	}
	gpnKeyId := ""
	if setting.SignWhenFirstCommit() {
		gpnKeyId = GetGpnKeyId(repo.Path, FirstCommitScene)
	}
	return CommitAndPushRepository(ctx, repo.Owner, Repository{
		Id:    repo.Id,
		Owner: repo.Owner,
		Name:  repo.Name,
		Path:  tmpDir,
	}, opts.DefaultBranch, "first commit", gpnKeyId)
}

func EnsureValidRepository(ctx context.Context, repoPath string) error {
	cmd := command.NewCommand("rev-parse")
	_, err := cmd.Run(ctx, command.WithDir(repoPath))
	return err
}

func SetDefaultBranch(ctx context.Context, repoPath, branch string) error {
	cmd := command.NewCommand("symbolic-ref", "HEAD", BranchPrefix+branch)
	_, err := cmd.Run(ctx, command.WithDir(repoPath))
	return err
}

func GetDefaultBranch(ctx context.Context, repoPath string) (string, error) {
	cmd := command.NewCommand("symbolic-ref", "HEAD")
	result, err := cmd.Run(ctx, command.WithDir(repoPath))
	if err != nil {
		return "", err
	}
	branch := result.ReadAsString()
	if !strings.HasPrefix(branch, BranchPrefix) {
		return "", fmt.Errorf("%s is not branch", branch)
	}
	return strings.TrimPrefix(strings.TrimSpace(branch), BranchPrefix), nil
}

func CloneRepository(ctx context.Context, repoPath, dst string, env []string) error {
	_, err := command.NewCommand("clone", repoPath, dst).Run(ctx, command.WithEnv(env))
	return err
}

func CommitAndPushRepository(ctx context.Context, committer User, repo Repository, branch, commitMsg string, gpnKeyId string) error {
	commitTimeStr := time.Now().Format(time.RFC3339)
	env := append(os.Environ(),
		"GIT_AUTHOR_NAME="+repo.Owner.Name,
		"GIT_AUTHOR_EMAIL="+repo.Owner.Email,
		"GIT_AUTHOR_DATE="+commitTimeStr,
		"GIT_COMMITTER_DATE="+commitTimeStr,
	)
	_, err := command.NewCommand("add", "--all").Run(ctx, command.WithDir(repo.Path))
	if err != nil {
		return fmt.Errorf("git add -all failed repo:%s err: %v", repo.Path, err)
	}
	commitCmd := command.NewCommand(
		"commit",
		fmt.Sprintf("--message='%s'", commitMsg),
		fmt.Sprintf("--author='%s <%s>'", repo.Owner.Name, repo.Owner.Email),
	)
	if gpnKeyId != "" {
		commitCmd.AddArgs("-S%s", gpnKeyId)
	} else {
		commitCmd.AddArgs("--no-gpg-sign")
	}
	env = append(env,
		"GIT_COMMITTER_NAME="+committer.Name,
		"GIT_COMMITTER_EMAIL="+committer.Email,
	)
	_, err = commitCmd.Run(ctx, command.WithDir(repo.Path), command.WithEnv(env))
	if err != nil {
		return fmt.Errorf("git commit failed repo:%s err: %v", repo.Path, err)
	}
	if branch == "" {
		branch = setting.DefaultBranch()
	}
	_, err = command.NewCommand("push", "origin", fmt.Sprintf("HEAD:%s", branch)).
		Run(ctx, command.WithDir(repo.Path), command.WithEnv(InternalPushEnv(repo, committer)))
	return err
}

func GetRepoUsername(repoPath string) (string, error) {
	return getRepoProperty("user.name", repoPath)
}

func GetRepoSignKey(repoPath string) (string, error) {
	return getRepoProperty("user.signingkey", repoPath)
}

func GetRepoUserEmail(repoPath string) (string, error) {
	return getRepoProperty("user.email", repoPath)
}

func getRepoProperty(name, repoPath string) (string, error) {
	run, err := command.NewCommand("config", "--get", name).Run(nil, command.WithDir(repoPath))
	if err != nil {
		return "", err
	}
	return run.ReadAsString(), nil
}

func FullPushEnv(repo Repository, committer User, prId string) []string {
	isWiki := "false"
	if strings.HasSuffix(repo.Path, ".wiki") {
		isWiki = "true"
	}
	environ := append(os.Environ(),
		"GIT_AUTHOR_NAME="+repo.Owner.Name,
		"GIT_AUTHOR_EMAIL="+repo.Owner.Email,
		"GIT_COMMITTER_NAME="+committer.Name,
		"GIT_COMMITTER_EMAIL="+committer.Email,
		EnvRepoName+"="+repo.Name,
		EnvRepoUsername+"="+repo.Owner.Name,
		EnvRepoIsWiki+"="+isWiki,
		EnvPusherName+"="+committer.Name,
		EnvPusherEmail+"="+committer.Email,
		EnvPusherID+"="+committer.Id,
		EnvRepoID+"="+repo.Id,
		EnvPRID+"="+prId,
		EnvAppURL+"="+setting.AppUrl(),
		"SSH_ORIGINAL_COMMAND=gitea-internal",
	)
	return environ
}

func InternalPushEnv(repo Repository, committer User) []string {
	return append(FullPushEnv(repo, committer, ""),
		EnvIsInternal+"=true",
	)
}
