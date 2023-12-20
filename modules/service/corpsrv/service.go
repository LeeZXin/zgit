package corpsrv

import (
	"context"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"zgit/modules/model/corpmd"
	"zgit/pkg/apicode"
	"zgit/pkg/i18n"
)

func GetCorpInfoByCorpId(ctx context.Context, id string) (corpmd.CorpInfo, bool, error) {
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	corp, b, err := corpmd.GetByCorpId(ctx, id)
	if err != nil {
		logger.Logger.WithContext(ctx).Errorf("GetCorpInfoByCorpId err: %v", err)
		return corpmd.CorpInfo{}, false, bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	if !b {
		return corpmd.CorpInfo{}, false, nil
	}
	return corp.ToCorpInfo(), true, nil
}
