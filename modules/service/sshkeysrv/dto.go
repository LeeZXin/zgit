package sshkeysrv

import (
	"github.com/LeeZXin/zsf-utils/bizerr"
	"zgit/modules/model/sshkeymd"
	"zgit/modules/model/usermd"
	"zgit/pkg/apicode"
	"zgit/pkg/i18n"
	"zgit/util"
)

type InsertSshKeyReqDTO struct {
	Name          string
	PubKeyContent string
	KeyType       int
	Operator      usermd.UserInfo
}

func (r *InsertSshKeyReqDTO) IsValid() error {
	if !util.AtLeastOneCharPattern.MatchString(r.Name) || len(r.Name) > 32 {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SshKeyInvalidName))
	}
	if r.PubKeyContent == "" {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SshKeyFormatError))
	}
	if r.KeyType != sshkeymd.UserPubKeyType && r.KeyType != sshkeymd.ProxyKeyType {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SshKeyInvalidKeyType))
	}
	if r.Operator.UserId == "" {
		return bizerr.NewBizErr(apicode.NotLoginCode.Int(), i18n.GetByKey(i18n.SystemNotLogin))
	}
	return nil
}

type DeleteSshKeyReqDTO struct {
	KeyId    string
	Operator usermd.UserInfo
}

func (r *DeleteSshKeyReqDTO) IsValid() error {
	if !util.AtLeastOneCharPattern.MatchString(r.KeyId) {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SshKeyInvalidKeyId))
	}
	if r.Operator.UserId == "" {
		return bizerr.NewBizErr(apicode.NotLoginCode.Int(), i18n.GetByKey(i18n.SystemNotLogin))
	}
	return nil
}
