package gitsrv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"github.com/gliderlabs/ssh"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kballard/go-shellquote"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"zgit/pkg/git"
	"zgit/pkg/git/command"
	"zgit/pkg/git/lfs"
	"zgit/pkg/git/process"
	"zgit/pkg/i18n"
	"zgit/pkg/perm"
	"zgit/setting"
	"zgit/standalone/modules/model/projectmd"
	"zgit/standalone/modules/model/repomd"
	"zgit/standalone/modules/model/usermd"
	"zgit/util"
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
		return errors.New(i18n.GetByKey(i18n.SystemInvalidArgs))
	}
	return after(ctx, keyUser, words, session)
}

func HandleGitCommand(ctx context.Context, operator usermd.UserInfo, words []string, session ssh.Session) error {
	verb := words[0]
	repoPath := strings.TrimPrefix(words[1], "/")
	var lfsVerb string
	if verb == lfsAuthenticateVerb {
		if !setting.LfsEnabled() {
			return errors.New(i18n.GetByKey(i18n.LfsNotSupported))
		}
		if len(words) > 2 {
			lfsVerb = words[2]
		}
	}
	accessMode, b := allowedCommands[verb]
	if !b {
		logger.Logger.Error("unsupported cmd: ", words)
		return errors.New(i18n.GetByKey(i18n.SshCmdNotSupported))
	}
	if verb == lfsAuthenticateVerb {
		if lfsVerb == "upload" {
			accessMode = perm.AccessModeWrite
		} else if lfsVerb == "download" {
			accessMode = perm.AccessModeRead
		} else {
			logger.Logger.Info("unsupported cmd: ", words)
			return errors.New(i18n.GetByKey(i18n.SshCmdNotSupported))
		}
	}
	repo, err := checkAccessMode(ctx, operator, repoPath, accessMode)
	if err != nil {
		return err
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
			RepoId:  repoPath,
			Op:      lfsVerb,
			Account: operator.Account,
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
			git.EnvRepoId, repo.RepoId,
			git.EnvPusherId, operator.Account,
			git.EnvAppUrl, setting.AppUrl(),
			git.EnvHookToken, setting.HookToken(),
		)...,
	)
	gitCmd.Env = append(gitCmd.Env, command.CommonEnvs()...)
	return gitCmd.Run()
}

func checkAccessMode(ctx context.Context, user usermd.UserInfo, repoPath string, accessMode perm.AccessMode) (repomd.Repo, error) {
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByPath(ctx, repoPath)
	if err != nil {
		logger.Logger.Error(err)
		return repomd.Repo{}, util.InternalError()
	}
	if !b {
		return repomd.Repo{}, util.InvalidArgsError()
	}
	// 系统管理员有所有的权限
	if user.IsAdmin {
		return repo, nil
	}
	p, b, err := projectmd.GetProjectUserPermDetail(ctx, repo.ProjectId, user.Account)
	if err != nil {
		logger.Logger.Error(err)
		return repomd.Repo{}, util.InternalError()
	}
	if !b {
		return repomd.Repo{}, util.UnauthorizedError()
	}
	if accessMode == perm.AccessModeWrite {
		// 检查权限
		if !p.PermDetail.GetRepoPerm(repo.RepoId).CanPush {
			return repomd.Repo{}, util.UnauthorizedError()
		}
	} else if accessMode == perm.AccessModeRead {
		// 检查权限
		if !p.PermDetail.GetRepoPerm(repo.RepoId).CanAccess {
			return repomd.Repo{}, util.UnauthorizedError()
		}
	} else {
		return repomd.Repo{}, util.InvalidArgsError()
	}
	return repo, nil
}
