package projectsrv

import (
	"context"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"github.com/patrickmn/go-cache"
	"time"
	"zgit/standalone/modules/model/projectmd"
	"zgit/util"
)

var (
	projectUserCache = cache.New(time.Minute, 10*time.Minute)
)

func InsertProject(ctx context.Context, reqDTO InsertProjectReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	err := projectmd.InsertProject(ctx, projectmd.InsertProjectReqDTO{
		Name: reqDTO.Name,
		Desc: reqDTO.Desc,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Errorf("insert project: %v", err)
		return util.InternalError()
	}
	return nil
}

func ProjectUserExists(ctx context.Context, projectId, account string) (bool, error) {
	cacheKey := projectId + "_" + account
	ret, b := projectUserCache.Get(cacheKey)
	if b {
		return ret.(bool), nil
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	exists, err := projectmd.ProjectUserExists(ctx, projectId, account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		projectUserCache.Set(cacheKey, false, time.Second)
		return false, err
	}
	projectUserCache.Set(cacheKey, exists, time.Minute)
	return exists, nil
}
