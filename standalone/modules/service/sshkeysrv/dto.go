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
	if !util.AtLeastOneCharPattern.MatchString(r.Name) || len(r.Name) > 32 {
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
	if !util.AtLeastOneCharPattern.MatchString(r.KeyId) {
		return util.InvalidArgsError()
	}
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	return nil
}
