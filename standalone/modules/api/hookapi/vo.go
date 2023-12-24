package hookapi

import "github.com/LeeZXin/zsf-utils/ginutil"

type RevInfo struct {
	OldCommitId string `json:"oldCommitId"`
	NewCommitId string `json:"newCommitId"`
	RefName     string `json:"refName"`
}

type OptsReqVO struct {
	RevInfoList []RevInfo `json:"revInfoList"`
	RepoId      string    `json:"repoId"`
	PrId        string    `json:"prId"`
	PusherId    string    `json:"pusherId"`
	IsWiki      bool      `json:"isWiki"`
}

type HttpRespVO struct {
	ginutil.BaseResp
}
