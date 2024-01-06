package sshkeymd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func GenKeyId() string {
	return idutil.RandomUuid()
}

func IsKeyIdValid(keyId string) bool {
	return len(keyId) == 32
}

func SearchByKeyContent(ctx context.Context, content string) (SshKey, bool, error) {
	var ret SshKey
	b, err := xormutil.MustGetXormSession(ctx).
		Where("content like ?", content+"%").
		Get(&ret)
	return ret, b, err
}

func GetByKeyId(ctx context.Context, keyId string) (SshKey, bool, error) {
	var ret SshKey
	b, err := xormutil.MustGetXormSession(ctx).
		Where("key_id = ?", keyId).
		Get(&ret)
	return ret, b, err
}

func DeleteSshKey(ctx context.Context, key SshKey) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("key_id = ?", key.KeyId).
		Delete(new(SshKey))
	return rows == 1, err
}

func InsertSshKey(ctx context.Context, reqDTO InsertSshKeyReqDTO) (SshKey, error) {
	p := SshKey{
		KeyId:       GenKeyId(),
		Account:     reqDTO.Account,
		Name:        reqDTO.Name,
		Fingerprint: reqDTO.Fingerprint,
		Content:     reqDTO.Content,
		Verified:    reqDTO.Verified,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&p)
	return p, err
}

func ListSshKey(ctx context.Context, reqDTO ListSshKeyReqDTO) ([]SshKey, error) {
	ret := make([]SshKey, 0)
	session := xormutil.MustGetXormSession(ctx).Where("account = ?", reqDTO.Account)
	if reqDTO.Offset > 0 {
		session.And("id > ?", reqDTO.Offset)
	}
	if reqDTO.Limit > 0 {
		session.Limit(reqDTO.Limit)
	}
	return ret, session.OrderBy("id asc").Find(&ret)
}

func UpdateVerifiedVar(ctx context.Context, reqDTO UpdateVerifiedVarReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("key_id = ?", reqDTO.KeyId).
		Cols("verified").
		Limit(1).
		Update(&SshKey{
			Verified: reqDTO.Verified,
		})
	return rows == 1, err
}
