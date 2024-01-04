package sshkeyapi

import "github.com/LeeZXin/zsf-utils/ginutil"

type InsertSshKeyReqVO struct {
	Name          string `json:"name"`
	PubKeyContent string `json:"pubKeyContent"`
}

type DeleteSshKeyReqVO struct {
	KeyId string `json:"keyId"`
}

type ListSshKeyReqVO struct {
	Offset int64 `json:"offset"`
	Limit  int   `json:"limit"`
}

type SshKeyVO struct {
	KeyId       string `json:"keyId"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
}

type ListSshKeyRespVO struct {
	ginutil.BaseResp
	Data   []SshKeyVO `json:"data"`
	Cursor int64      `json:"cursor"`
}

type GetTokenReqVO struct {
	KeyId string `json:"keyId"`
}

type GetTokenRespVO struct {
	ginutil.BaseResp
	Token string `json:"token"`
}

type VerifyTokenReqVO struct {
	KeyId     string `json:"keyId"`
	Token     string `json:"token"`
	Signature string `json:"signature"`
}
