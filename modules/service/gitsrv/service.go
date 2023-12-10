package gitsrv

import (
	"context"
	"errors"
	"fmt"
	"github.com/gliderlabs/ssh"
	"github.com/kballard/go-shellquote"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"zgit/modules/model/usermd"
	"zgit/pkg/git"
	"zgit/pkg/git/command"
	"zgit/pkg/git/process"
	"zgit/pkg/perm"
	"zgit/setting"
	"zgit/util"
)

const (
	lfsAuthenticateVerb = "git-lfs-authenticate"
)

var (
	hiWords = "Hi there! You've successfully authenticated with the deploy key named %v, but zgit does not provide shell access."

	allowedCommands = map[string]perm.AccessMode{
		"git-upload-pack":    perm.AccessModeRead,
		"git-upload-archive": perm.AccessModeRead,
		"git-receive-pack":   perm.AccessModeWrite,
		lfsAuthenticateVerb:  perm.AccessModeNone,
	}
)

func HandleSshCommand(ctx context.Context, cmd string, keyUser usermd.UserInfo, session ssh.Session, after func(context.Context, usermd.UserInfo, []string, ssh.Session) error) error {
	// 命令为空
	if cmd == "" {
		fmt.Fprintln(session, fmt.Sprintf(hiWords, keyUser.Name))
		return nil
	}
	words, err := shellquote.Split(cmd)
	if err != nil {
		return errors.New("error parsing arguments")
	}
	if len(words) < 2 {
		if git.CheckGitVersionAtLeast("2.29") == nil {
			if cmd == "ssh_info" {
				fmt.Fprintln(session, `{"type":"zgit","version":1}`)
				return nil
			}
		}
		return errors.New("too few arguments")
	}
	return after(ctx, keyUser, words, session)
}

func HandleGitCommand(ctx context.Context, operator usermd.UserInfo, words []string, session ssh.Session) error {
	verb := words[0]
	repoPath := strings.TrimPrefix(words[1], "/")
	var lfsVerb string
	if verb == lfsAuthenticateVerb {
		if !setting.LfsEnabled() {
			return errors.New("LFS authentication request over SSH denied, LFS support is disabled")
		}
		if len(words) > 2 {
			lfsVerb = words[2]
		}
	}
	relativeRepoPath, err := util.ParseRelativeRepoPath(repoPath)
	if err != nil {
		return errors.New("Invalid repository path:" + repoPath)
	}
	accessMode, b := allowedCommands[verb]
	if !b {
		return errors.New("Unsupported git command:" + verb)
	}
	if verb == lfsAuthenticateVerb {
		if lfsVerb == "upload" {
			accessMode = perm.AccessModeWrite
		} else if lfsVerb == "download" {
			accessMode = perm.AccessModeRead
		} else {
			return errors.New("Unknown LFS verb:" + lfsVerb)
		}
	}
	results, b, err := checkAccessMode(ctx, operator, relativeRepoPath, accessMode, verb, lfsVerb)
	if err != nil {
		return errors.New("not authorized")
	}
	var gitcmd *exec.Cmd
	gitBinPath := filepath.Dir(setting.GitExecutablePath()) // e.g. /usr/bin
	gitBinVerb := filepath.Join(gitBinPath, verb)           // e.g. /usr/bin/git-upload-pack
	if _, err = os.Stat(gitBinVerb); err != nil {
		verbFields := strings.SplitN(verb, "-", 2)
		if len(verbFields) == 2 {
			gitcmd = exec.CommandContext(ctx, setting.GitExecutablePath(), verbFields[1], repoPath)
		}
	}
	if gitcmd == nil {
		gitcmd = exec.CommandContext(ctx, gitBinVerb, repoPath)
	}
	process.SetSysProcAttribute(gitcmd)
	gitcmd.Dir = setting.RepoDir()
	gitcmd.Stdout = session
	gitcmd.Stdin = session
	gitcmd.Stderr = session.Stderr()
	gitcmd.Env = append(gitcmd.Env, os.Environ()...)
	gitcmd.Env = append(gitcmd.Env,
		util.JoinFields(
			git.EnvRepoIsWiki, strconv.FormatBool(results.IsWiki),
			git.EnvRepoPath, filepath.Join(setting.RepoDir(), repoPath),
			git.EnvRepoID, results.RepoId,
			git.EnvPusherID, operator.Id,
		)...,
	)
	gitcmd.Env = append(gitcmd.Env, command.CommonEnvs()...)
	return gitcmd.Run()
}

func checkAccessMode(ctx context.Context, user usermd.UserInfo, repo util.RelativeRepoPath, accessMode perm.AccessMode, verbs ...string) (ServCommandResults, bool, error) {
	return ServCommandResults{}, true, nil
}

type ServCommandResults struct {
	IsWiki    bool
	RepoId    string
	ClusterId string
}
