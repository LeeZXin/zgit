package proxy

import (
	"github.com/LeeZXin/zsf-utils/quit"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	gossh "golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
	"zgit/setting"
	"zgit/util"
)

var (
	serverCiphers      = []string{"chacha20-poly1305@openssh.com", "aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "aes256-gcm@openssh.com"}
	serverKeyExchanges = []string{"curve25519-sha256", "ecdh-sha2-nistp256", "ecdh-sha2-nistp384", "ecdh-sha2-nistp521", "diffie-hellman-group14-sha256", "diffie-hellman-group14-sha1"}
	serverMACs         = []string{"hmac-sha2-256-etm@openssh.com", "hmac-sha2-256", "hmac-sha1"}
	serverHostKey      = "ssh/proxy.rsa"
	serverPort         = static.GetInt("ssh.proxy.port")
	proxyName          = static.GetString("ssh.proxy.name")

	serverKeySigner gossh.Signer

	clientConfig *gossh.ClientConfig
)

func StartSSHProxy() {
	if serverPort <= 0 {
		logger.Logger.Panic("ssh proxy port should greater than 0")
	}
	if !filepath.IsAbs(serverHostKey) {
		serverHostKey = filepath.Join(setting.DataDir(), serverHostKey)
	}
	if err := os.MkdirAll(filepath.Dir(serverHostKey), os.ModePerm); err != nil {
		logger.Logger.Panicf("Failed to create dir %s: %v", serverHostKey, err)
	}
	exist, err := util.IsExist(serverHostKey)
	if err != nil {
		logger.Logger.Panicf("check host key failed %s: %v", serverHostKey, err)
	}
	if !exist {
		err = util.GenKeyPair(serverHostKey)
		if err != nil {
			logger.Logger.Panicf("gen host key pair failed %s: %v", serverHostKey, err)
		}
	}
	privateKey, err := os.ReadFile(serverHostKey)
	if err != nil {
		logger.Logger.Panicf("read pub key failed %s: %v", serverHostKey, err)
	}
	serverKeySigner, err = gossh.ParsePrivateKey(privateKey)
	if err != nil {
		logger.Logger.Panicf("parse signer failed %s: %v", serverHostKey, err)
	}
	clientConfig = &gossh.ClientConfig{
		Config: gossh.Config{
			KeyExchanges: serverKeyExchanges,
			Ciphers:      serverCiphers,
			MACs:         serverMACs,
		},
		User: "git",
		Auth: []gossh.AuthMethod{
			gossh.PublicKeys(serverKeySigner),
		},
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
	}
	s := newProxy()
	quit.AddShutdownHook(s.Shutdown)
	s.Start()
	quit.Wait()
}
