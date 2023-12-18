package gitsrv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LeeZXin/zsf/logger"
	"github.com/gliderlabs/ssh"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kballard/go-shellquote"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"zgit/modules/model/usermd"
	"zgit/pkg/git"
	"zgit/pkg/git/command"
	"zgit/pkg/git/process"
	"zgit/pkg/lfs"
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
	logger.Logger.Info("git cmd: ", words)
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
	// LFS token authentication
	if verb == lfsAuthenticateVerb {
		url := fmt.Sprintf("%s/%s/info/lfs", setting.AppUrl(), repoPath)
		now := time.Now()
		claims := lfs.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(setting.LfsJwtAuthExpiry())),
				NotBefore: jwt.NewNumericDate(now),
			},
			RepoId: results.RepoId,
			Op:     lfsVerb,
			UserId: operator.Id,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Sign and get the complete encoded token as a string using the secret
		tokenStr, err := token.SignedString(setting.LfsJwtSecretBytes())
		if err != nil {
			return fmt.Errorf("failed to sign JWT token: %v", err)
		}
		tokenAuthentication := &lfs.TokenRespVO{
			Header: map[string]string{
				"Authorization": tokenStr,
			},
			Href: url,
		}
		err = json.NewEncoder(session).Encode(tokenAuthentication)
		if err != nil {
			return fmt.Errorf("failed to encode LFS json response: %v", err)
		}
		return nil
	}
	var gitCmd *exec.Cmd
	gitBinPath := filepath.Dir(setting.GitExecutablePath()) // e.g. /usr/bin
	gitBinVerb := filepath.Join(gitBinPath, verb)           // e.g. /usr/bin/git-upload-pack
	if _, err = os.Stat(gitBinVerb); err != nil {
		verbFields := strings.SplitN(verb, "-", 2)
		if len(verbFields) == 2 {
			gitCmd = exec.CommandContext(ctx, setting.GitExecutablePath(), verbFields[1], repoPath)
		}
	}
	if gitCmd == nil {
		gitCmd = exec.CommandContext(ctx, gitBinVerb, repoPath)
	}
	process.SetSysProcAttribute(gitCmd)
	gitCmd.Dir = setting.RepoDir()
	gitCmd.Stdout = session
	gitCmd.Stdin = session
	gitCmd.Stderr = session.Stderr()
	gitCmd.Env = append(gitCmd.Env, os.Environ()...)
	gitCmd.Env = append(gitCmd.Env,
		util.JoinFields(
			git.EnvRepoIsWiki, strconv.FormatBool(results.IsWiki),
			git.EnvRepoPath, filepath.Join(setting.RepoDir(), repoPath),
			git.EnvRepoID, results.RepoId,
			git.EnvPusherID, operator.Id,
			git.EnvAppUrl, setting.AppUrl(),
		)...,
	)
	gitCmd.Env = append(gitCmd.Env, command.CommonEnvs()...)
	return gitCmd.Run()
}

func checkAccessMode(ctx context.Context, user usermd.UserInfo, repo util.RelativeRepoPath, accessMode perm.AccessMode, verbs ...string) (ServCommandResults, bool, error) {
	return ServCommandResults{}, true, nil
}

type ServCommandResults struct {
	IsWiki    bool
	RepoId    string
	ClusterId string
}
