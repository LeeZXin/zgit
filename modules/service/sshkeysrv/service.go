package sshkeysrv

import (
	"context"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	gossh "golang.org/x/crypto/ssh"
	"strings"
	"zgit/modules/model/sshkeymd"
	"zgit/pkg/apicode"
	"zgit/pkg/i18n"
)

func SearchByKeyContent(ctx context.Context, key gossh.PublicKey, keyType int) (sshkeymd.SshKey, bool, error) {
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	pubKey, b, err := sshkeymd.SearchByKeyTypeAndContent(ctx, keyType, strings.TrimSpace(string(gossh.MarshalAuthorizedKey(key))))
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return sshkeymd.SshKey{}, false, bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	return pubKey, b, nil
}

func DeleteSshKey(ctx context.Context, reqDTO DeleteSshKeyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	sshKey, b, err := sshkeymd.GetByKeyId(ctx, reqDTO.KeyId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	if !b {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SshKeyNotFound))
	}
	// 只有拥有人才能删掉公钥
	if sshKey.UserId != reqDTO.Operator.UserId {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SystemUnauthorized))
	}
	_, err = sshkeymd.DeleteSshKey(ctx, sshKey)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	return nil
}

func InsertSshKey(ctx context.Context, reqDTO InsertSshKeyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	publicKey, _, _, _, err := gossh.ParseAuthorizedKey([]byte(reqDTO.PubKeyContent))
	if err != nil {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SshKeyFormatError))
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	_, b, err := SearchByKeyContent(ctx, publicKey, reqDTO.KeyType)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	if b {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SshKeyAlreadyExists))
	}
	fingerprint := gossh.FingerprintSHA256(publicKey)
	_, err = sshkeymd.InsertSshKey(ctx, sshkeymd.InsertSshKeyReqDTO{
		UserId:      reqDTO.Operator.UserId,
		CorpId:      reqDTO.Operator.CorpId,
		Name:        reqDTO.Name,
		Fingerprint: fingerprint,
		Content:     strings.TrimSpace(string(gossh.MarshalAuthorizedKey(publicKey))),
		KeyType:     reqDTO.KeyType,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	return nil
}
