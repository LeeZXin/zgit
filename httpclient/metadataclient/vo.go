package metadataclient

import "github.com/LeeZXin/zsf-utils/ginutil"

type GetCorpInfoReqVO struct {
	CorpId string `json:"corpId"`
}

type CorpInfoVO struct {
	CorpId     string `json:"corpId"`
	Name       string `json:"name"`
	NodeId     string `json:"nodeId"`
	RepoCount  int    `json:"repoCount"`
	RepoLimit  int    `json:"repoLimit"`
	MaxLfsSize int    `json:"maxLfsSize"`
	MaxGitSize int    `json:"maxGitSize"`
}

type GetCorpInfoRespVO struct {
	ginutil.BaseResp
	IsExists bool       `json:"isExists"`
	Corp     CorpInfoVO `json:"corp"`
}

type GetClusterNodeInfoReqVO struct {
	NodeId string `json:"nodeId"`
}

type ClusterNodeInfoVO struct {
	NodeId string `json:"nodeId"`
	Name   string `json:"name"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
}

type GetClusterNodeInfoRespVO struct {
	ginutil.BaseResp
	IsExists bool              `json:"isExists"`
	Node     ClusterNodeInfoVO `json:"node"`
}
