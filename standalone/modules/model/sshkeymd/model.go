package sshkeymd

import (
	"time"
)

const (
	SshKeyTableName = "ssh_key"
)

type KeyInfo struct {
	KeyId       string `json:"keyId"`
	Account     string `json:"account"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
	Content     string `json:"content"`
	Verified    bool   `json:"verified"`
}

type SshKey struct {
	Id          int64     `xorm:"pk autoincr"`
	KeyId       string    `json:"keyId"`
	Account     string    `json:"account"`
	Name        string    `json:"name"`
	Fingerprint string    `json:"fingerprint"`
	Content     string    `json:"content"`
	Verified    bool      `json:"verified"`
	Created     time.Time `json:"created" xorm:"created"`
	Updated     time.Time `json:"updated" xorm:"updated"`
}

func (k *SshKey) ToKeyInfo() KeyInfo {
	return KeyInfo{
		KeyId:       k.KeyId,
		Account:     k.Account,
		Name:        k.Name,
		Fingerprint: k.Fingerprint,
		Content:     k.Content,
		Verified:    k.Verified,
	}
}

func (*SshKey) TableName() string {
	return SshKeyTableName
}
