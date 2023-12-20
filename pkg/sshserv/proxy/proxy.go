package proxy

import (
	"context"
	"errors"
	"fmt"
	"github.com/LeeZXin/zsf/logger"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"zgit/modules/model/clustermd"
	"zgit/modules/model/sshkeymd"
	"zgit/modules/model/usermd"
	"zgit/modules/service/gitsrv"
	"zgit/modules/service/reposrv"
	"zgit/modules/service/sshkeysrv"
	"zgit/modules/service/usersrv"
	"zgit/pkg/sshserv"
	"zgit/setting"
	"zgit/util"
)

type NodeInfo struct {
	Id   string `json:"id"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

func publicKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	if ctx.User() != setting.GitUser() {
		return false
	}
	pubKey, b, err := sshkeysrv.SearchByKeyContent(ctx, key, sshkeymd.UserPubKeyType)
	if !b || err != nil {
		return false
	}
	userInfo, b, err := usersrv.GetUserInfoByUserId(ctx, pubKey.UserId)
	if !b || err != nil {
		return false
	}
	ctx.SetValue(sshserv.ZgitUserId, userInfo)
	return true
}

func sessionHandler(session ssh.Session) {
	ctx := session.Context()
	userInfo := session.Context().Value(sshserv.ZgitUserId).(usermd.UserInfo)
	if err := gitsrv.HandleSshCommand(ctx, session.RawCommand(), userInfo, session, handleProxyCommand); err != nil {
		util.ExitWithErrMsg(session, err.Error())
	} else {
		session.Exit(0)
	}
}

func handleProxyCommand(ctx context.Context, operator usermd.UserInfo, words []string, session ssh.Session) error {
	repoPath := strings.TrimPrefix(words[1], "/")
	repoInfo, b, err := reposrv.GetRepoInfoByPath(ctx, repoPath)
	if !b || err != nil {
		return fmt.Errorf("could not find repo: %s", repoPath)
	}
	nodeInfo, b, err := clustermd.GetByNodeId(ctx, repoInfo.NodeId)
	if !b || err != nil {
		return fmt.Errorf("could not nodeInfo: %s", repoPath)
	}
	// 建立SSH连接
	client, err := gossh.Dial("tcp", fmt.Sprintf("%s:%d", nodeInfo.Host, nodeInfo.Port), clientConfig)
	if err != nil {
		return errors.New("connect to proxy failed")
	}
	defer client.Close()
	proxySession, err := client.NewSession()
	if err != nil {
		return errors.New("connect to proxy failed")
	}
	defer proxySession.Close()
	for _, env := range session.Environ() {
		b, a, f := strings.Cut(env, "=")
		if f {
			if err = proxySession.Setenv(b, a); err != nil {
				return fmt.Errorf("can not transfer env name: %s", b)
			}
		}
	}
	if err = proxySession.Setenv("ZGIT_PROXY_NAME", proxyName); err != nil {
		return errors.New("can not transfer proxy name:" + proxyName)
	}
	if err = proxySession.Setenv("ZGIT_LOGIN_USER", operator.UserId); err != nil {
		return errors.New("can not transfer login user")
	}
	stdout, err := proxySession.StdoutPipe()
	if err != nil {
		return errors.New("network err")
	}
	stderr, err := proxySession.StderrPipe()
	if err != nil {
		return errors.New("network err")
	}
	stdin, err := proxySession.StdinPipe()
	if err != nil {
		return errors.New("network err")
	}
	defer stdin.Close()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	if err = proxySession.Start(session.RawCommand()); err != nil {
		return errors.New("network err")
	}
	go func() {
		defer stdin.Close()
		io.Copy(stdin, session)
	}()
	go func() {
		defer wg.Done()
		io.Copy(session, stdout)
	}()
	go func() {
		defer wg.Done()
		io.Copy(session.Stderr(), stderr)
	}()
	wg.Wait()
	return proxySession.Wait()
}

type proxy struct {
	*ssh.Server
}

func newProxy() *proxy {
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
	return &proxy{
		Server: srv,
	}
}

func (s *proxy) Start() {
	go func() {
		logger.Logger.Infof("start ssh proxy port: %d", serverPort)
		err := s.ListenAndServe()
		if err != nil && err != ssh.ErrServerClosed {
			logger.Logger.Panicf("ssh proxy err: %v", err)
		}
	}()
}

func (s *proxy) Shutdown() {
	logger.Logger.Info("shutdown ssh proxy")
	s.Close()
}
