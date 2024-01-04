package sshkeysrv

import (
	"strings"
	"zgit/pkg/git/signature"
	"zgit/standalone/modules/model/usermd"
	"zgit/util"
)

type InsertSshKeyReqDTO struct {
	Name          string
	PubKeyContent string
	Operator      usermd.UserInfo
}

func (r *InsertSshKeyReqDTO) IsValid() error {
	if len(r.Name) == 0 || len(r.Name) > 128 {
		return util.InvalidArgsError()
	}
	if r.PubKeyContent == "" {
		return util.InvalidArgsError()
	}
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteSshKeyReqDTO struct {
	KeyId    string
	Operator usermd.UserInfo
}

func (r *DeleteSshKeyReqDTO) IsValid() error {
	if !validateKeyId(r.KeyId) {
		return util.InvalidArgsError()
	}
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type ListSshKeyReqDTO struct {
	Offset   int64
	Limit    int
	Operator usermd.UserInfo
}

func (r *ListSshKeyReqDTO) IsValid() error {
	if r.Offset < 0 {
		return util.InvalidArgsError()
	}
	if r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type ListSshKeyRespDTO struct {
	Cursor  int64
	KeyList []SshKeyDTO
}

type SshKeyDTO struct {
	KeyId       string
	Name        string
	Fingerprint string
}

type GetTokenReqDTO struct {
	KeyId    string
	Operator usermd.UserInfo
}

func (r *GetTokenReqDTO) IsValid() error {
	if !validateKeyId(r.KeyId) {
		return util.InvalidArgsError()
	}
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type VerifySshKeyReqDTO struct {
	KeyId     string
	Token     string
	Signature string
	Operator  usermd.UserInfo
}

func (r *VerifySshKeyReqDTO) IsValid() error {
	if !validateKeyId(r.KeyId) {
		return util.InvalidArgsError()
	}
	if len(r.Token) != 72 {
		return util.InvalidArgsError()
	}
	r.Signature = strings.TrimSpace(r.Signature)
	if !strings.HasPrefix(r.Signature, signature.StartSSHSigLineTag) || !strings.HasSuffix(r.Signature, signature.EndSSHSigLineTag) {
		return util.InvalidArgsError()
	}
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

func validateKeyId(keyId string) bool {
	return len(keyId) == 32
}

func validateOperator(operator usermd.UserInfo) bool {
	return operator.Account != ""
}
