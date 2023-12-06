package proxy

import (
	"fmt"
	"github.com/LeeZXin/zsf/logger"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"zgit/setting"
)

type contextKey string

const (
	zgitKeyId   = contextKey("zgit-key-id")
	zgitCluster = contextKey("zgit-proxy-cluster")
)

type UserInfo struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	ClusterId string `json:"clusterId"`
}

type NodeInfo struct {
	Id   string `json:"id"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

func sshConnectionFailed(net.Conn, error) {}

func publicKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	if ctx.User() != setting.GitUser() {
		return false
	}
	keyContent := strings.TrimSpace(string(gossh.MarshalAuthorizedKey(key)))
	userInfo, err := getUserInfoByPublicKey(keyContent)
	if err != nil {
		return false
	}
	nodeInfo, err := getNodeInfoByClusterId(userInfo.ClusterId)
	if err != nil {
		return false
	}
	ctx.SetValue(zgitCluster, nodeInfo)
	ctx.SetValue(zgitKeyId, userInfo.Id)
	return true
}

func getUserInfoByPublicKey(pubKey string) (*UserInfo, error) {
	return &UserInfo{
		Name:  "zexin",
		Email: "zexin@fake.local",
	}, nil
}

func getNodeInfoByClusterId(clusterId string) (*NodeInfo, error) {
	return &NodeInfo{
		Id:   "1",
		Host: "127.0.0.1",
		Port: 3333,
	}, nil
}

func sessionHandler(session ssh.Session) {
	ctx := session.Context()
	nodeInfo := ctx.Value(zgitCluster).(*NodeInfo)
	// 建立SSH连接
	client, err := gossh.Dial("tcp", fmt.Sprintf("%s:%d", nodeInfo.Host, nodeInfo.Port), clientConfig)
	if err != nil {
		exitWithErrMsg(session, "connect to proxy failed")
		return
	}
	defer client.Close()
	proxySession, err := client.NewSession()
	if err != nil {
		exitWithErrMsg(session, "connect to proxy failed")
		return
	}
	defer proxySession.Close()
	for _, env := range session.Environ() {
		b, a, f := strings.Cut(env, "=")
		if f {
			proxySession.Setenv(b, a)
		}
	}
	proxySession.Setenv("ZGIT_PROXY_NAME", proxyName)
	proxySession.Setenv("ZGIT_LOGIN_USER", ctx.Value(zgitKeyId).(string))
	stdout, err := proxySession.StdoutPipe()
	if err != nil {
		exitWithErrMsg(session, "network err")
		return
	}
	stderr, err := proxySession.StderrPipe()
	if err != nil {
		exitWithErrMsg(session, "network err")
		return
	}
	stdin, err := proxySession.StdinPipe()
	if err != nil {
		exitWithErrMsg(session, "network err")
		return
	}
	defer stdin.Close()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	if err = proxySession.Start(session.RawCommand()); err != nil {
		exitWithErrMsg(session, "network err")
		return
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
	if err = proxySession.Wait(); err != nil {
		exitWithErrMsg(session, err.Error())
	} else {
		session.Exit(0)
	}
}

func exitWithErrMsg(session ssh.Session, msg string) {
	fmt.Fprintln(session.Stderr(), msg)
	session.Exit(1)
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

func (s *server) Start() {
	go func() {
		logger.Logger.Infof("start ssh server port: %d", serverPort)
		err := s.ListenAndServe()
		if err != nil && err != ssh.ErrServerClosed {
			logger.Logger.Panicf("ssh server err: %v", err)
		}
	}()
}

func (s *server) Shutdown() {
	logger.Logger.Info("shutdown ssh server")
	s.Close()
}
