package sshkeymd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
)

func SearchByKeyTypeAndContent(ctx context.Context, keyType int, content string) (SshKey, bool, error) {
	var ret SshKey
	b, err := mysqlstore.GetXormSession(ctx).
		Where("content like ?", content+"%").
		And("key_type = ?", keyType).
		Get(&ret)
	return ret, b, err
}

func GetByKeyId(ctx context.Context, keyId string) (SshKey, bool, error) {
	var ret SshKey
	b, err := mysqlstore.GetXormSession(ctx).
		Where("key_id = ?", keyId).
		Get(&ret)
	return ret, b, err
}

func DeleteSshKey(ctx context.Context, key SshKey) (bool, error) {
	rows, err := mysqlstore.GetXormSession(ctx).
		Where("key_id = ?", key.KeyId).
		Limit(1).
		Delete(new(SshKey))
	return rows == 1, err
}

func InsertSshKey(ctx context.Context, reqDTO InsertSshKeyReqDTO) (SshKey, error) {
	p := SshKey{
		KeyId:       GenKeyId(),
		UserId:      reqDTO.UserId,
		CorpId:      reqDTO.CorpId,
		Name:        reqDTO.Name,
		Fingerprint: reqDTO.Fingerprint,
		Content:     reqDTO.Content,
		KeyType:     reqDTO.KeyType,
		Verified:    reqDTO.Verified,
	}
	_, err := mysqlstore.GetXormSession(ctx).Insert(&p)
	return p, err
}
