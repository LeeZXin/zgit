package sshserv

import (
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/zsf"
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
		err = util.GenKeyPair(serverHostKey)
		if err != nil {
			logger.Logger.Panicf("gen host key pair failed %s: %v", serverHostKey, err)
		}
	}
	zsf.RegisterApplicationLifeCycle(newServer())
}
