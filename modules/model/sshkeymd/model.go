package sshkeymd

import (
	"github.com/LeeZXin/zsf-utils/idutil"
	"time"
)

const (
	SshKeyTableName = "ssh_key"
)

const (
	UserPubKeyType = iota
	ProxyKeyType
)

func GenKeyId() string {
	return idutil.RandomUuid()
}

type SshKey struct {
	Id          int64     `xorm:"pk autoincr"`
	KeyId       string    `json:"keyId"`
	UserId      string    `json:"userId"`
	CorpId      string    `json:"corpId"`
	Name        string    `json:"name"`
	Fingerprint string    `json:"fingerprint"`
	Content     string    `json:"content"`
	KeyType     int       `json:"keyType"`
	Verified    bool      `json:"verified"`
	Created     time.Time `json:"created" xorm:"created"`
	Updated     time.Time `json:"updated" xorm:"updated"`
}

func (*SshKey) TableName() string {
	return SshKeyTableName
}
