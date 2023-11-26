package ssh

import (
	"context"
	"errors"
	"fmt"
	"github.com/LeeZXin/zsf/logger"
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

const zgitKeyID = contextKey("zgit-key-id")

func sshConnectionFailed(net.Conn, error) {}

func publicKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	isDebugRunMode := setting.IsDebugRunMode()
	if isDebugRunMode {
		logger.Logger.Debugf("Handle Public Key: Fingerprint: %s from %s with user %s", gossh.FingerprintSHA256(key), ctx.RemoteAddr(), ctx.User())
	}
	if ctx.User() != setting.GitUser() {
		return false
	}
	if isDebugRunMode {
		logger.Logger.Debug(strings.TrimSpace(string(gossh.MarshalAuthorizedKey(key))))
	}
	ctx.SetValue(zgitKeyID, int64(1))
	return true
}

func sessionHandler(session ssh.Session) {
	isDebugRunMode := setting.IsDebugRunMode()
	keyID := fmt.Sprintf("%d", session.Context().Value(zgitKeyID).(int64))
	command := session.RawCommand()
	if isDebugRunMode {
		logger.Logger.Debug("SSH: Payload: %v", command)
	}
	logger.Logger.Infof("SSH: Payload: %v", command)
	args := []string{"serv", "key-" + keyID}
	ctx, cancel := context.WithCancel(session.Context())
	defer cancel()
	gitProtocol := ""
	for _, env := range session.Environ() {
		if strings.HasPrefix(env, "GIT_PROTOCOL=") {
			_, gitProtocol, _ = strings.Cut(env, "=")
			break
		}
	}
	cmd := exec.CommandContext(ctx, setting.AppPath(), args...)
	cmd.Env = append(
		os.Environ(),
		"SSH_ORIGINAL_COMMAND="+command,
		"SKIP_MINWINSVC=1",
		"GIT_PROTOCOL="+gitProtocol,
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Logger.Error("SSH: StdoutPipe: %v", err)
		return
	}
	defer stdout.Close()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Logger.Error("SSH: StderrPipe: %v", err)
		return
	}
	defer stderr.Close()
	stdin, err := cmd.StdinPipe()
	if err != nil {
		logger.Logger.Error("SSH: StdinPipe: %v", err)
		return
	}
	defer stdin.Close()
	process.SetSysProcAttribute(cmd)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	if isDebugRunMode {
		logger.Logger.Debugf("cmd: %v", cmd.String())
	}
	if err = cmd.Start(); err != nil {
		logger.Logger.Error("SSH: Start: %v", err)
		return
	}
	go func() {
		defer stdin.Close()
		if _, err := io.Copy(stdin, session); err != nil {
			logger.Logger.Error("Failed to write session to stdin. %s", err)
		}
	}()
	go func() {
		defer wg.Done()
		defer stdout.Close()
		if _, err := io.Copy(session, stdout); err != nil {
			logger.Logger.Error("Failed to write stdout to session. %s", err)
		}
	}()
	go func() {
		defer wg.Done()
		defer stderr.Close()
		if _, err := io.Copy(session.Stderr(), stderr); err != nil {
			logger.Logger.Error("Failed to write stderr to session. %s", err)
		}
	}()
	wg.Wait()
	err = cmd.Wait()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			logger.Logger.Error("SSH: Wait: %v", err)
		}
	}
	if err = session.Exit(getExitStatusFromError(err)); err != nil && !errors.Is(err, io.EOF) {
		logger.Logger.Error("Session failed to exit. %s", err)
	}
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
