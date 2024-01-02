package cfgsrv

import (
	"context"
	"errors"
	"github.com/LeeZXin/zsf-utils/localcache"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"time"
	"zgit/standalone/modules/model/syscfgmd"
	"zgit/util"
)

var (
	sysCfgCache *localcache.SingleCacheEntry[*SysCfg]
)

func init() {
	sysCfgCache, _ = localcache.NewSingleCacheEntry(func(ctx context.Context) (*SysCfg, error) {
		cfg, b, err := GetSysCfgFromDB(ctx)
		if !b || err != nil {
			return nil, errors.New("sysCfg is not existed")
		}
		return &cfg, nil
	}, 30*time.Second)
}

func InitSysCfg() {
	ctx := context.Background()
	cfg, b, err := GetSysCfgFromDB(ctx)
	if err != nil {
		logger.Logger.Panic(err)
	}
	if !b {
		cfg = SysCfg{}
		err = InsertSysCfg(ctx, cfg)
		if err != nil {
			logger.Logger.Panic(err)
		}
	}
}

func InsertSysCfg(ctx context.Context, cfg SysCfg) error {
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	err := syscfgmd.InsertCfg(ctx, &cfg)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	return nil
}

func GetSysCfgWithCache(ctx context.Context) (*SysCfg, error) {
	return sysCfgCache.LoadData(ctx)
}

func GetSysCfgFromDB(ctx context.Context) (SysCfg, bool, error) {
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	ret := SysCfg{}
	b, err := syscfgmd.GetByKey(ctx, &ret)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return SysCfg{}, false, util.InternalError()
	}
	return ret, b, err
}

func GetSysCfg(ctx context.Context, reqDTO GetSysCfgReqDTO) (SysCfg, error) {
	if err := reqDTO.IsValid(); err != nil {
		return SysCfg{}, err
	}
	if !reqDTO.Operator.IsAdmin {
		return SysCfg{}, util.UnauthorizedError()
	}
	cfg, _, err := GetSysCfgFromDB(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return SysCfg{}, err
	}
	return cfg, nil
}

func UpdateSysCfg(ctx context.Context, reqDTO UpdateSysCfgReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	_, err := syscfgmd.UpdateByKey(ctx, &reqDTO.SysCfg)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return err
	}
	sysCfgCache.Clear()
	return nil
}
