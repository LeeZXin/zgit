package syscfgmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"zgit/util"
)

func InsertCfg(ctx context.Context, kv util.KeyVal) error {
	ret := SysCfg{
		CfgKey:  kv.Key(),
		Content: kv.Val(),
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return err
}

func GetByKey(ctx context.Context, kv util.KeyVal) (bool, error) {
	ret := SysCfg{}
	b, err := xormutil.MustGetXormSession(ctx).Where("cfg_key = ?", kv.Key()).Get(&ret)
	if b {
		err = kv.FromStore(ret.Content)
	}
	return b, err
}

func UpdateByKey(ctx context.Context, kv util.KeyVal) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("cfg_key = ?", kv.Key()).
		Cols("content").
		Limit(1).
		Update(&SysCfg{
			Content: kv.Val(),
		})
	return rows == 1, err
}

func DeleteByKey(ctx context.Context, key string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("cfg_key = ?", key).Limit(1).Delete(new(SysCfg))
	return rows == 1, err
}
