package clustersrv

import (
	"context"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"zgit/modules/model/clustermd"
	"zgit/pkg/apicode"
	"zgit/pkg/i18n"
)

func GetClusterInfoById(ctx context.Context, id string) (clustermd.NodeInfo, bool, error) {
	node, b, err := clustermd.GetByNodeId(ctx, id)
	if err != nil {
		return clustermd.NodeInfo{}, false, bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	return node, b, nil
}
