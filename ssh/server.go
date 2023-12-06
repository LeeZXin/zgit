package ssh

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"zgit/git/process"
	"zgit/setting"
)

type contextKey string

const (
	zgitKeyId = contextKey("zgit-key-id")

	standaloneMode = "standalone"
)

var (
	mode = static.GetString("cluster.mode")
)

func init() {
	if mode == "" {
		mode = standaloneMode
	}
}

type UserInfo struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	ClusterId string `json:"clusterId"`
}

func sshConnectionFailed(net.Conn, error) {}

func publicKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	if ctx.User() != setting.GitUser() {
		return false
	}
	keyContent := strings.TrimSpace(string(gossh.MarshalAuthorizedKey(key)))
	if mode == standaloneMode {
		userInfo, b, err := getUserInfoByPublicKey(keyContent)
		if !b || err != nil {
			return false
		}
		ctx.SetValue(zgitKeyId, userInfo.Id)
	} else {
		b, err := existNodeInfoByPublicKey(keyContent)
		if !b || err != nil {
			return false
		}
	}
	return true
}

func getUserInfoByPublicKey(pubKey string) (*UserInfo, bool, error) {
	return &UserInfo{
		Name:  "zexin",
		Email: "zexin@fake.local",
	}, true, nil
}

func existNodeInfoByPublicKey(pubKey string) (bool, error) {
	return true, nil
}

func sessionHandler(session ssh.Session) {
	ctx, cancel := context.WithCancel(session.Context())
	defer cancel()
	var keyID string
	if mode == standaloneMode {
		keyID = session.Context().Value(zgitKeyId).(string)
	} else {
		for _, env := range session.Environ() {
			if strings.HasPrefix(env, "ZGIT_LOGIN_USER") {
				_, after, f := strings.Cut(env, "ZGIT_LOGIN_USER")
				if f {
					keyID = after
				}
			}
		}
	}
	if keyID == "" {
		fmt.Fprintf(session.Stderr(), "lost login user")
		session.Exit(1)
		return
	}
	command := session.RawCommand()
	cmd := exec.CommandContext(ctx, setting.AppPath(), "serv", "key-"+keyID)
	cmd.Env = append(
		os.Environ(),
		"SSH_ORIGINAL_COMMAND="+command,
		"SKIP_MINWINSVC=1",
	)
	cmd.Env = append(
		cmd.Env,
		session.Environ()...,
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	defer stdout.Close()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}
	defer stderr.Close()
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	defer stdin.Close()
	process.SetSysProcAttribute(cmd)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	if err = cmd.Start(); err != nil {
		return
	}
	go func() {
		defer stdin.Close()
		io.Copy(stdin, session)
	}()
	go func() {
		defer wg.Done()
		defer stdout.Close()
		io.Copy(session, stdout)
	}()
	go func() {
		defer wg.Done()
		defer stderr.Close()
		io.Copy(session.Stderr(), stderr)
	}()
	wg.Wait()
	err = cmd.Wait()
	session.Exit(getExitStatusFromError(err))
}

func getExitStatusFromError(err error) int {
	if err == nil {
		return 0
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return 1
	}
	waitStatus, ok := exitErr.Sys().(syscall.WaitStatus)
	if !ok {
		if exitErr.Success() {
			return 0
		}
		return 1
	}
	return waitStatus.ExitStatus()
}

type server struct {
	*ssh.Server
}

func newServer() *server {
	srv := &ssh.Server{
		Addr:             net.JoinHostPort("", strconv.Itoa(serverPort)),
		PublicKeyHandler: publicKeyHandler,
		Handler:          sessionHandler,
		ServerConfigCallback: func(ctx ssh.Context) *gossh.ServerConfig {
			config := &gossh.ServerConfig{}
			config.KeyExchanges = serverKeyExchanges
			config.MACs = serverMACs
			config.Ciphers = serverCiphers
			return config
		},
		ConnectionFailedCallback: sshConnectionFailed,
		PtyCallback: func(ctx ssh.Context, pty ssh.Pty) bool {
			return false
		},
	}
	if err := srv.SetOption(ssh.HostKeyFile(serverHostKey)); err != nil {
		logger.Logger.Panic(err)
	}
	return &server{
		Server: srv,
	}
}

func (s *server) OnApplicationStart() {
	go func() {
		logger.Logger.Infof("start ssh server port: %d", serverPort)
		err := s.ListenAndServe()
		if err != nil && err != ssh.ErrServerClosed {
			logger.Logger.Panicf("ssh server err: %v", err)
		}
	}()
}

func (s *server) AfterInitialize() {}

func (s *server) OnApplicationShutdown() {
	logger.Logger.Info("shutdown ssh server")
	s.Close()
}
