package metadataclient

import (
	"context"
	"github.com/LeeZXin/zsf/http/httpclient"
	"github.com/LeeZXin/zsf/property/static"
)

var (
	client        = httpclient.Dial("metadata-http")
	authorization = httpclient.WithHeader(map[string]string{
		"Authorization": static.GetString("httpclient.metadata.token"),
	})
)

// GetCorpInfo 获取企业信息
func GetCorpInfo(ctx context.Context, reqVO GetCorpInfoReqVO) (respVO GetCorpInfoRespVO, err error) {
	err = client.Post(ctx, "/metadata/corp/get", reqVO, &respVO, authorization)
	return
}

// GetClusterNodeInfo 获取集群节点信息
func GetClusterNodeInfo(ctx context.Context, reqVO GetClusterNodeInfoReqVO) (respVO GetClusterNodeInfoRespVO, err error) {
	err = client.Post(ctx, "/metadata/cluster/get", reqVO, &respVO, authorization)
	return
}
