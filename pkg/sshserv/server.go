package sshserv

import (
	"context"
	"github.com/LeeZXin/zsf/logger"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"net"
	"strconv"
	"strings"
	"zgit/modules/model/sshkeymd"
	"zgit/modules/model/usermd"
	"zgit/modules/service/gitsrv"
	"zgit/modules/service/sshkeysrv"
	"zgit/modules/service/usersrv"
	"zgit/setting"
	"zgit/util"
)

type ContextKey string

const (
	ZgitUserId = ContextKey("zgit-user-id")

	standaloneMode = "standalone"
	shardingMode   = "sharding"
)

var (
	mode string
)

func publicKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	if ctx.User() != setting.GitUser() {
		return false
	}
	if mode == standaloneMode {
		pubKey, b, err := sshkeysrv.SearchByKeyContent(ctx, key, sshkeymd.UserPubKeyType)
		if !b || err != nil {
			return false
		}
		userInfo, b, err := usersrv.GetUserInfoByUserId(ctx, pubKey.UserId)
		if !b || err != nil {
			return false
		}
		ctx.SetValue(ZgitUserId, userInfo)
	} else {
		_, b, err := sshkeysrv.SearchByKeyContent(ctx, key, sshkeymd.ProxyKeyType)
		if !b || err != nil {
			return false
		}
	}
	return true
}

func sessionHandler(session ssh.Session) {
	ctx, cancel := context.WithCancel(session.Context())
	defer cancel()
	var userInfo usermd.UserInfo
	if mode == standaloneMode {
		userInfo = session.Context().Value(ZgitUserId).(usermd.UserInfo)
	} else {
		var userId string
		for _, env := range session.Environ() {
			if strings.HasPrefix(env, "ZGIT_LOGIN_USER") {
				_, after, f := strings.Cut(env, "ZGIT_LOGIN_USER")
				if f {
					userId = after
				}
			}
		}
		if userId == "" {
			util.ExitWithErrMsg(session, "lost login user\n")
			return
		}
		info, b, err := usersrv.GetUserInfoByUserId(ctx, userId)
		if !b || err != nil {
			if err != nil {
				logger.Logger.Errorf("GetUserInfoByUserId err: %v", err)
			}
			util.ExitWithErrMsg(session, "lost login user\n")
			return
		}
		userInfo = info
	}
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
