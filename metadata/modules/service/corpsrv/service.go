package corpsrv

import (
	"context"
	"github.com/LeeZXin/zsf/logger"
	"github.com/patrickmn/go-cache"
	"time"
	"zgit/httpclient/metadataclient"
	"zgit/util"
)

var (
	corpInfoCache = cache.New(time.Minute, 10*time.Minute)
)

func GetCorpInfoByCorpId(ctx context.Context, id string) (CorpInfoDTO, bool, error) {
	v, b := corpInfoCache.Get(id)
	if b {
		ret := v.(CorpInfoDTO)
		// 判断是否是空缓存
		if ret.CorpId == "" {
			return ret, false, nil
		}
		return ret, true, nil
	}
	respVO, err := metadataclient.GetCorpInfo(ctx, metadataclient.GetCorpInfoReqVO{
		CorpId: id,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return CorpInfoDTO{}, false, util.InternalError()
	}
	if respVO.IsExists {
		ret := CorpInfoDTO{
			CorpId:     respVO.Corp.CorpId,
			Name:       respVO.Corp.Name,
			NodeId:     respVO.Corp.NodeId,
			RepoCount:  respVO.Corp.RepoCount,
			RepoLimit:  respVO.Corp.RepoLimit,
			MaxLfsSize: respVO.Corp.MaxLfsSize,
			MaxGitSize: respVO.Corp.MaxGitSize,
		}
		corpInfoCache.Set(id, ret, 3*time.Minute)
		return ret, true, nil
	}
	return CorpInfoDTO{}, false, nil
}
