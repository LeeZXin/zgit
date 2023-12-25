package sshkeysrv

import (
	"context"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"github.com/patrickmn/go-cache"
	gossh "golang.org/x/crypto/ssh"
	"strings"
	"time"
	"zgit/pkg/apicode"
	"zgit/pkg/i18n"
	"zgit/standalone/modules/model/sshkeymd"
)

var (
	sshKeyCache = cache.New(time.Minute, 10*time.Minute)
)

func SearchByKeyContent(ctx context.Context, key gossh.PublicKey) (sshkeymd.SshKey, bool, error) {
	keyContent := strings.TrimSpace(string(gossh.MarshalAuthorizedKey(key)))
	sshKey, b := sshKeyCache.Get(keyContent)
	if b {
		ret := sshKey.(sshkeymd.SshKey)
		if ret.Id == 0 {
			return ret, false, nil
		}
		return ret, true, nil
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	pubKey, b, err := sshkeymd.SearchByKeyContent(ctx, keyContent)
	if err != nil {
		// 空缓存
		sshKeyCache.Set(keyContent, sshkeymd.SshKey{}, time.Second)
		logger.Logger.WithContext(ctx).Error(err)
		return sshkeymd.SshKey{}, false, bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	sshKeyCache.Set(keyContent, pubKey, 3*time.Minute)
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
	if sshKey.Account != reqDTO.Operator.Account {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SystemUnauthorized))
	}
	_, err = sshkeymd.DeleteSshKey(ctx, sshKey)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	// 删除缓存
	sshKeyCache.Delete(sshKey.Content)
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
	_, b, err := SearchByKeyContent(ctx, publicKey)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	if b {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SshKeyAlreadyExists))
	}
	fingerprint := gossh.FingerprintSHA256(publicKey)
	_, err = sshkeymd.InsertSshKey(ctx, sshkeymd.InsertSshKeyReqDTO{
		Account:     reqDTO.Operator.Account,
		Name:        reqDTO.Name,
		Fingerprint: fingerprint,
		Content:     strings.TrimSpace(string(gossh.MarshalAuthorizedKey(publicKey))),
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	return nil
}
