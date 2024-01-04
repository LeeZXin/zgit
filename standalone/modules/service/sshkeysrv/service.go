package sshkeysrv

import (
	"context"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	gossh "golang.org/x/crypto/ssh"
	"strings"
	"time"
	"zgit/pkg/apicode"
	"zgit/pkg/i18n"
	"zgit/pkg/signature"
	"zgit/standalone/modules/model/sshkeymd"
	"zgit/util"
)

var (
	sshKeyCache = util.NewGoCache()
	tokenCache  = util.NewGoCache()
)

func SearchByKeyContent(ctx context.Context, key gossh.PublicKey) (sshkeymd.KeyInfo, bool, error) {
	keyContent := strings.TrimSpace(string(gossh.MarshalAuthorizedKey(key)))
	sshKey, b := sshKeyCache.Get(keyContent)
	if b {
		ret := sshKey.(sshkeymd.KeyInfo)
		if ret.KeyId == "" {
			return ret, false, nil
		}
		return ret, true, nil
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	pubKey, b, err := sshkeymd.SearchByKeyContent(ctx, keyContent)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return sshkeymd.KeyInfo{}, false, bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	if !b {
		// 空缓存
		k := sshkeymd.KeyInfo{}
		sshKeyCache.Set(keyContent, k, time.Second)
		return k, false, nil
	}
	ret := pubKey.ToKeyInfo()
	sshKeyCache.Set(keyContent, ret, time.Minute)
	return ret, b, nil
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
		return util.InternalError()
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 只有拥有人才能删掉公钥
	if sshKey.Account != reqDTO.Operator.Account {
		return util.InvalidArgsError()
	}
	_, err = sshkeymd.DeleteSshKey(ctx, sshKey)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
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
		return util.InternalError()
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
		return util.InternalError()
	}
	return nil
}

func ListSshKey(ctx context.Context, reqDTO ListSshKeyReqDTO) (ListSshKeyRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return ListSshKeyRespDTO{}, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 展示登录人的sshkey列表
	keyList, err := sshkeymd.ListSshKey(ctx, sshkeymd.ListSshKeyReqDTO{
		Offset:  reqDTO.Offset,
		Limit:   reqDTO.Limit,
		Account: reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return ListSshKeyRespDTO{}, util.InternalError()
	}
	ret := ListSshKeyRespDTO{}
	ret.KeyList, _ = listutil.Map(keyList, func(t sshkeymd.SshKey) (SshKeyDTO, error) {
		return SshKeyDTO{
			KeyId:       t.KeyId,
			Name:        t.Name,
			Fingerprint: t.Fingerprint,
		}, nil
	})
	if len(keyList) > 0 {
		ret.Cursor = keyList[len(keyList)-1].Id
	}
	return ret, nil
}

func GetToken(ctx context.Context, reqDTO GetTokenReqDTO) (string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	sshKey, b, err := sshkeymd.GetByKeyId(ctx, reqDTO.KeyId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError()
	}
	if !b {
		return "", util.InvalidArgsError()
	}
	if sshKey.Account != reqDTO.Operator.Account {
		return "", util.InvalidArgsError()
	}
	token := signature.GetToken(signature.User{
		Account: reqDTO.Operator.Account,
		Email:   reqDTO.Operator.Email,
	})
	// 设置十分钟有效期
	tokenCache.Set(reqDTO.KeyId, token, 10*time.Minute)
	return token, nil
}

// VerifySshKey 校验ssh key
func VerifySshKey(ctx context.Context, reqDTO VerifySshKeyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	sshKey, b, err := sshkeymd.GetByKeyId(ctx, reqDTO.KeyId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if !b {
		return util.InvalidArgsError()
	}
	if sshKey.Account != reqDTO.Operator.Account {
		return util.InvalidArgsError()
	}
	// 已经校验过了
	if sshKey.Verified {
		return bizerr.NewBizErr(apicode.SshKeyAlreadyVerifiedCode.Int(), i18n.GetByKey(i18n.SshKeyAlreadyVerified))
	}
	//首先校验token正确
	if !signature.VerifyToken(reqDTO.Token, signature.User{
		Account: reqDTO.Operator.Account,
		Email:   reqDTO.Operator.Email,
	}) {
		return util.InvalidArgsError()
	}
	_, b = tokenCache.Get(reqDTO.KeyId)
	// token不存在或已失效
	if !b {
		return bizerr.NewBizErr(apicode.SshKeyVerifyTokenExpiredCode.Int(), i18n.GetByKey(i18n.SshKeyVerifyTokenExpired))
	}
	// 校验签名
	if err = signature.VerifySshSignature(reqDTO.Signature, reqDTO.Token, sshKey.Content); err != nil {
		// 校验失败
		return bizerr.NewBizErr(apicode.SshKeyVerifyFailedCode.Int(), i18n.GetByKey(i18n.SshKeyVerifyFailed))
	}
	if _, err = sshkeymd.UpdateVerifiedVar(ctx, sshkeymd.UpdateVerifiedVarReqDTO{
		KeyId:    reqDTO.KeyId,
		Verified: true,
	}); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	return nil
}
