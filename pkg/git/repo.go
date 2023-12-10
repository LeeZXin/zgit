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
	"zgit/pkg/git/command"
	"zgit/setting"
	"zgit/util"
)

var (
	gitIgnoreResourcesPath = filepath.Join(setting.ResourcesDir(), "gitignore")
)

type Repository struct {
	Id    string `json:"Id"`
	Owner User   `json:"owner"`
	Name  string `json:"name"`
	Path  string `json:"path"`
}

type InitRepoOpts struct {
	Owner         User
	RepoName      string
	RepoPath      string
	CreateReadme  bool
	GitIgnoreName string
	DefaultBranch string
}

type CommitAndPushOpts struct {
	RepoPath          string
	Owner, Committer  User
	Branch, CommitMsg string
}

func initEmptyRepository(ctx context.Context, repoPath string, bare bool) error {
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
	cmd := command.NewCommand("init")
	if bare {
		cmd.AddArgs(command.BareFlag)
	}
	_, err = cmd.Run(ctx, command.WithDir(repoPath))
	return err
}

func InitRepository(ctx context.Context, opts InitRepoOpts) error {
	if err := initEmptyRepository(ctx, opts.RepoPath, true); err != nil {
		return err
	}
	if opts.CreateReadme || opts.GitIgnoreName != "" {
		tmpDir, err := os.MkdirTemp(setting.TempDir(), "zgit-"+opts.RepoName)
		if err != nil {
			return fmt.Errorf("failed to create temp dir for repository %s: %w", opts.RepoPath, err)
		}
		defer util.RemoveAll(tmpDir)
		if err = initTemporaryRepository(ctx, tmpDir, opts); err != nil {
			return err
		}
	} else {
		branch := opts.DefaultBranch
		if branch == "" {
			branch = setting.DefaultBranch()
		}
		SetDefaultBranch(ctx, opts.RepoPath, branch)
	}
	return InitRepoHook(opts.RepoPath)
}

func initTemporaryRepository(ctx context.Context, tmpDir string, opts InitRepoOpts) error {
	if _, err := command.NewCommand("clone", opts.RepoPath, tmpDir).Run(ctx); err != nil {
		return fmt.Errorf("failed to clone original repository %s: %w", opts.RepoPath, err)
	}
	if opts.CreateReadme {
		util.WriteFile(filepath.Join(tmpDir, "README.md"), []byte(fmt.Sprintf("# %s  ", opts.RepoName)))
	}
	if opts.GitIgnoreName != "" {
		content, err := os.ReadFile(filepath.Join(gitIgnoreResourcesPath, opts.GitIgnoreName))
		if err == nil {
			util.WriteFile(filepath.Join(tmpDir, ".gitIgnore"), content)
		}
	}
	return commitAndPushRepository(ctx, CommitAndPushOpts{
		RepoPath:  tmpDir,
		Owner:     opts.Owner,
		Committer: opts.Owner,
		Branch:    opts.DefaultBranch,
		CommitMsg: "first commit",
	})
}

func EnsureValidRepository(ctx context.Context, repoPath string) error {
	cmd := command.NewCommand("rev-parse")
	_, err := cmd.Run(ctx, command.WithDir(repoPath))
	return err
}

func SetDefaultBranch(ctx context.Context, repoPath, branch string) error {
	if !strings.HasPrefix(branch, BranchPrefix) {
		branch = BranchPrefix + branch
	}
	cmd := command.NewCommand("symbolic-ref", "HEAD", branch)
	_, err := cmd.Run(ctx, command.WithDir(repoPath))
	return err
}

func commitAndPushRepository(ctx context.Context, opts CommitAndPushOpts) error {
	commitTimeStr := time.Now().Format(time.RFC3339)
	env := append(
		os.Environ(),
		util.JoinFields(
			"GIT_AUTHOR_NAME", opts.Owner.Name,
			"GIT_AUTHOR_EMAIL", opts.Owner.Email,
			"GIT_AUTHOR_DATE", commitTimeStr,
			"GIT_COMMITTER_DATE", commitTimeStr,
		)...,
	)
	_, err := command.NewCommand("add", "--all").Run(ctx, command.WithDir(opts.RepoPath))
	if err != nil {
		return fmt.Errorf("git add -all failed repo:%s err: %v", opts.RepoPath, err)
	}
	commitCmd := command.NewCommand(
		"commit",
		"--no-gpg-sign",
		fmt.Sprintf("--message='%s'", opts.CommitMsg),
		fmt.Sprintf("--author='%s <%s>'", opts.Committer.Name, opts.Committer.Email),
	)
	env = append(env,
		util.JoinFields(
			"GIT_COMMITTER_NAME", opts.Committer.Name,
			"GIT_COMMITTER_EMAIL", opts.Committer.Email,
		)...,
	)
	_, err = commitCmd.Run(ctx, command.WithDir(opts.RepoPath), command.WithEnv(env))
	if err != nil {
		return fmt.Errorf("git commit failed repo:%s err: %v", opts.RepoPath, err)
	}
	if opts.Branch == "" {
		opts.Branch = setting.DefaultBranch()
	}
	_, err = command.NewCommand("push", "origin", opts.Branch).
		Run(ctx, command.WithDir(opts.RepoPath), command.WithEnv(util.JoinFields(EnvIsInternal, "true")))
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
