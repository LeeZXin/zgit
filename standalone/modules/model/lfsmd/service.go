package lfsmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func InsertLock(ctx context.Context, reqDTO InsertLockReqDTO) (LfsLock, error) {
	ret := LfsLock{
		RepoId: reqDTO.RepoId,
		Owner:  reqDTO.Owner,
		Path:   reqDTO.Path,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func GetLockById(ctx context.Context, id int64) (LfsLock, bool, error) {
	var ret LfsLock
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func DeleteLock(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Delete(new(LfsLock))
	return rows == 1, err
}
