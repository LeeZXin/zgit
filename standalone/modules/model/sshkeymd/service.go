package sshkeymd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

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
		Limit(1).
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
