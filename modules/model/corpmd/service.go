package corpmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
)

func InsertCorp(ctx context.Context, reqDTO InsertCorpReqDTO) (Corp, error) {
	c := Corp{
		CorpId:    reqDTO.CorpId,
		Name:      reqDTO.Name,
		CorpDesc:  reqDTO.CorpDesc,
		RepoLimit: reqDTO.RepoLimit,
	}
	_, err := mysqlstore.GetXormSession(ctx).Insert(&c)
	return c, err
}

func DeleteCorp(ctx context.Context, corpId string) error {
	_, err := mysqlstore.GetXormSession(ctx).Where("corp_id = ?", corpId).Limit(1).Delete(new(Corp))
	return err
}

func GetByCorpId(ctx context.Context, corpId string) (Corp, bool, error) {
	var ret Corp
	b, err := mysqlstore.GetXormSession(ctx).Where("corp_id = ?", corpId).Get(&ret)
	return ret, b, err
}

func IncrRepoCount(ctx context.Context, corpId string) (bool, error) {
	rows, err := mysqlstore.GetXormSession(ctx).Where("corp_id = ?", corpId).Incr("repo_count").Update(new(Corp))
	return rows == 1, err
}

func DecrRepoCount(ctx context.Context, corpId string) (bool, error) {
	rows, err := mysqlstore.GetXormSession(ctx).Where("corp_id = ?", corpId).Decr("repo_count").Update(new(Corp))
	return rows == 1, err
}
