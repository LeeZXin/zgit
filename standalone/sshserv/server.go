package sshserv

import (
	"context"
	"github.com/LeeZXin/zsf/logger"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"net"
	"strconv"
	"zgit/setting"
	"zgit/standalone/modules/model/usermd"
	"zgit/standalone/modules/service/gitsrv"
	"zgit/standalone/modules/service/sshkeysrv"
	"zgit/standalone/modules/service/usersrv"
	"zgit/util"
)

type ContextKey string

const (
	ZgitUserAccount = ContextKey("zgit-user-account")
)

func publicKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	if ctx.User() != setting.GitUser() {
		return false
	}
	pubKey, b, err := sshkeysrv.SearchByKeyContent(ctx, key)
	if !b || err != nil {
		return false
	}
	userInfo, b, err := usersrv.GetUserInfoByAccount(ctx, pubKey.Account)
	if !b || err != nil {
		return false
	}
	ctx.SetValue(ZgitUserAccount, userInfo)
	return true
}

func sessionHandler(session ssh.Session) {
	ctx, cancel := context.WithCancel(session.Context())
	defer cancel()
	userInfo := session.Context().Value(ZgitUserAccount).(usermd.UserInfo)
	if err := gitsrv.HandleSshCommand(ctx, session.RawCommand(), userInfo, session, gitsrv.HandleGitCommand); err != nil {
		util.ExitWithErrMsg(session, err.Error())
	} else {
		session.Exit(0)
	}
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
