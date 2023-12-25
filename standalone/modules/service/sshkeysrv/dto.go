package sshkeysrv

import (
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
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteSshKeyReqDTO struct {
	KeyId    string
	Operator usermd.UserInfo
}

func (r *DeleteSshKeyReqDTO) IsValid() error {
	if len(r.KeyId) == 0 || len(r.KeyId) > 32 {
		return util.InvalidArgsError()
	}
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	return nil
}
