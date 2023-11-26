package cmd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/LeeZXin/zsf/logger"
	"github.com/kballard/go-shellquote"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"zgit/git"
	"zgit/git/command"
	"zgit/git/process"
	"zgit/perm"
	"zgit/setting"
)

const (
	lfsAuthenticateVerb = "git-lfs-authenticate"
)

var (
	hiWords = "Hi there! You've successfully authenticated with the deploy key named %v, but zgit does not provide shell access."

	alphaDashDotPattern = regexp.MustCompile(`[^\w-\.]`)

	allowedCommands = map[string]perm.AccessMode{
		"git-upload-pack":    perm.AccessModeRead,
		"git-upload-archive": perm.AccessModeRead,
		"git-receive-pack":   perm.AccessModeWrite,
		lfsAuthenticateVerb:  perm.AccessModeNone,
	}
)

var Serv = &cli.Command{
	Name:        "serv",
	Usage:       "This command should only be called by SSH shell",
	Description: "Serv provides access auth for repositories",
	Action:      runServ,
}

func runServ(c *cli.Context) error {
	isDebugRunMode := setting.IsDebugRunMode()
	ctx, cancel := initWaitContext()
	defer cancel()
	if c.NArg() < 1 {
		return exitWithDefaultCode("too few arguments")
	}
	keys := strings.Split(c.Args().First(), "-")
	if len(keys) != 2 || keys[0] != "key" {
		return exitWithDefaultCode("Key ID format error")
	}
	keyUser, b, err := findAuthorizedKey(ctx, keys[1])
	if err != nil {
		return exitWithDefaultCode("Key ID format error")
	}
	if !b {
		return exitWithDefaultCode("Key check failed")
	}
	cmd := os.Getenv("SSH_ORIGINAL_COMMAND")
	if isDebugRunMode {
		logger.Logger.Debugf("runServ command: %v", cmd)
	}
	// 命令为空
	if cmd == "" {
		fmt.Printf(hiWords, keyUser.UserName)
		return nil
	}
	words, err := shellquote.Split(cmd)
	if err != nil {
		return exitWithDefaultCode("Error parsing arguments")
	}
	if len(words) < 2 {
		return sshInfo(cmd)
	}
	return gitAction(ctx, keyUser, isDebugRunMode, words)
}

func findAuthorizedKey(ctx context.Context, keyId string) (KeyUser, bool, error) {
	return KeyUser{
		UserName: "zxjcli",
	}, true, nil
}

type KeyUser struct {
	UserName string `json:"userName"`
}

func sshInfo(cmd string) error {
	if git.CheckGitVersionAtLeast("2.29") == nil {
		if cmd == "ssh_info" {
			fmt.Print(`{"type":"zgit","version":1}`)
			return nil
		}
	}
	return exitWithDefaultCode("too few arguments")
}

func gitAction(ctx context.Context, user KeyUser, isDebugMode bool, words []string) error {
	verb := words[0]
	repoPath := words[1]
	if repoPath[0] == '/' {
		repoPath = repoPath[1:]
	}
	var lfsVerb string
	if verb == lfsAuthenticateVerb {
		if !setting.LfsEnabled() {
			return exitWithDefaultCode("LFS authentication request over SSH denied, LFS support is disabled")
		}
		if len(words) > 2 {
			lfsVerb = words[2]
		}
	}
	repoPath = strings.ToLower(strings.TrimSpace(repoPath))
	rr := strings.SplitN(repoPath, "/", 2)
	if len(rr) != 2 {
		return exitWithDefaultCode("Invalid repository path" + repoPath)
	}
	username := strings.ToLower(rr[0])
	reponame := strings.ToLower(strings.TrimSuffix(rr[1], ".git"))
	if isDebugMode {
		logger.Logger.Debugf("repo: %s username: %s reponame: %s", repoPath, username, reponame)
	}
	if alphaDashDotPattern.MatchString(reponame) {
		return exitWithDefaultCode("Invalid repo name")
	}
	accessMode, b := allowedCommands[verb]
	if !b {
		return exitWithDefaultCode("Unsupported git command:" + verb)
	}
	if verb == lfsAuthenticateVerb {
		if lfsVerb == "upload" {
			accessMode = perm.AccessModeWrite
		} else if lfsVerb == "download" {
			accessMode = perm.AccessModeRead
		} else {
			return exitWithDefaultCode("Unknown LFS verb:" + lfsVerb)
		}
	}
	results, b, err := checkAccessMode(ctx, user, username, reponame, accessMode, verb, lfsVerb)
	if err != nil {
		return exitWithDefaultCode("not authorized")
	}
	// LFS token authentication
	//if verb == lfsAuthenticateVerb {
	//	url := fmt.Sprintf("%s%s/%s.git/info/lfs", setting.AppURL, url.PathEscape(results.OwnerName), url.PathEscape(results.Name))
	//	now := time.Now()
	//	claims := lfs.Claims{
	//		RegisteredClaims: jwt.RegisteredClaims{
	//			ExpiresAt: jwt.NewNumericDate(now.Add(setting.LFS.HTTPAuthExpiry)),
	//			NotBefore: jwt.NewNumericDate(now),
	//		},
	//		RepoID: results.RepoID,
	//		Op:     lfsVerb,
	//		UserID: results.UserID,
	//	}
	//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//
	//	// Sign and get the complete encoded token as a string using the secret
	//	tokenString, err := token.SignedString(setting.LFS.JWTSecretBytes)
	//	if err != nil {
	//		return fail(ctx, "Failed to sign JWT Token", "Failed to sign JWT token: %v", err)
	//	}
	//	tokenAuthentication := &git_model.LFSTokenResponse{
	//		Header: make(map[string]string),
	//		Href:   url,
	//	}
	//	tokenAuthentication.Header["Authorization"] = fmt.Sprintf("Bearer %s", tokenString)
	//	enc := json.NewEncoder(os.Stdout)
	//	err = enc.Encode(tokenAuthentication)
	//	if err != nil {
	//		return fail(ctx, "Failed to encode LFS json response", "Failed to encode LFS json response: %v", err)
	//	}
	//	return nil
	//}
	var gitcmd *exec.Cmd
	gitBinPath := filepath.Dir(setting.GitExecutablePath()) // e.g. /usr/bin
	gitBinVerb := filepath.Join(gitBinPath, verb)           // e.g. /usr/bin/git-upload-pack
	if _, err = os.Stat(gitBinVerb); err != nil {
		// if the command "git-upload-pack" doesn't exist, try to split "git-upload-pack" to use the sub-command with git
		// ps: Windows only has "git.exe" in the bin path, so Windows always uses this way
		verbFields := strings.SplitN(verb, "-", 2)
		if len(verbFields) == 2 {
			gitcmd = exec.CommandContext(ctx, setting.GitExecutablePath(), verbFields[1], repoPath)
		}
	}
	if gitcmd == nil {
		// by default, use the verb (it has been checked above by allowedCommands)
		gitcmd = exec.CommandContext(ctx, gitBinVerb, repoPath)
	}
	process.SetSysProcAttribute(gitcmd)
	stderr := &bytes.Buffer{}
	gitcmd.Dir = setting.RepoRootDir()
	gitcmd.Stdout = os.Stdout
	gitcmd.Stdin = os.Stdin
	gitcmd.Stderr = stderr
	gitcmd.Env = append(gitcmd.Env, os.Environ()...)
	gitcmd.Env = append(gitcmd.Env,
		git.EnvRepoIsWiki+"="+strconv.FormatBool(results.IsWiki),
		git.EnvRepoName+"="+results.RepoName,
		git.EnvRepoUsername+"="+results.OwnerName,
		git.EnvPusherName+"="+results.UserName,
		git.EnvPusherEmail+"="+results.UserEmail,
		git.EnvPusherID+"="+strconv.FormatInt(results.UserID, 10),
		git.EnvRepoID+"="+strconv.FormatInt(results.RepoID, 10),
		git.EnvPRID+"="+fmt.Sprintf("%d", 0),
		git.EnvDeployKeyID+"="+fmt.Sprintf("%d", results.DeployKeyID),
		git.EnvKeyID+"="+fmt.Sprintf("%d", results.KeyID),
		git.EnvAppURL+"="+setting.AppUrl(),
	)
	// to avoid breaking, here only use the minimal environment variables for the "gitea serv" command.
	// it could be re-considered whether to use the same git.CommonGitCmdEnvs() as "git" command later.
	gitcmd.Env = append(gitcmd.Env, command.CommonEnvs()...)
	if isDebugMode {
		logger.Logger.Debugf("gitcmd: %s", gitcmd.String())
	}
	if err = gitcmd.Run(); err != nil {
		return exitWithDefaultCode(fmt.Sprintf("Failed to execute git command: %v", stderr.String()))
	}
	return nil
}

func checkAccessMode(ctx context.Context, user KeyUser, username, reponame string, accessMode perm.AccessMode, verbs ...string) (ServCommandResults, bool, error) {
	return ServCommandResults{}, true, nil
}

type ServCommandResults struct {
	IsWiki      bool
	DeployKeyID int64
	KeyID       int64  // public key
	KeyName     string // this field is ambiguous, it can be the name of DeployKey, or the name of the PublicKey
	UserName    string
	UserEmail   string
	UserID      int64
	OwnerName   string
	RepoName    string
	RepoID      int64
}
