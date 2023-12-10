package clustersrv

import (
	"context"
	"zgit/modules/model/clustermd"
)

func GetClusterInfoById(ctx context.Context, id string) (clustermd.ClusterInfo, bool, error) {
	return clustermd.ClusterInfo{
		Id:   "1",
		Host: "127.0.0.1",
		Port: 3333,
	}, true, nil
}
