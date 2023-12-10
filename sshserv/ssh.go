package sshserv

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/zsf"
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
	serverHostKey      = "ssh/zgit.rsa"
	serverPort         = static.GetInt("ssh.server.port")
)

func InitSsh() {
	mode = static.GetString("cluster.mode")
	if mode == "" {
		mode = standaloneMode
	}
	if mode != standaloneMode && mode != shardingMode {
		logger.Logger.Panicf("unknown cluster mode: %s", mode)
	}
	logger.Logger.Infof("working in cluster mode: %s", mode)
	if serverPort <= 0 {
		logger.Logger.Panic("ssh server port should greater than 0")
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
		err = genKeyPair(serverHostKey)
		if err != nil {
			logger.Logger.Panicf("gen host key pair failed %s: %v", serverHostKey, err)
		}
	}
	zsf.RegisterApplicationLifeCycle(newServer())
}

func genKeyPair(keyPath string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	f, err := os.OpenFile(keyPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = pem.Encode(f, privateKeyPEM); err != nil {
		return err
	}
	pub, err := gossh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	public := gossh.MarshalAuthorizedKey(pub)
	p, err := os.OpenFile(keyPath+".pub", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer p.Close()
	_, err = p.Write(public)
	return err
}
