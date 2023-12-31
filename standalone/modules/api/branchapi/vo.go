package branchapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"zgit/standalone/modules/model/branchmd"
)

type InsertProtectedBranchReqVO struct {
	RepoId string                      `json:"repoId"`
	Branch string                      `json:"branch"`
	Cfg    branchmd.ProtectedBranchCfg `json:"cfg"`
}

type DeleteProtectedBranchReqVO struct {
	Bid string `json:"bid"`
}

type ListProtectedBranchReqVO struct {
	RepoId string `json:"repoId"`
}

type ProtectedBranchVO struct {
	Bid    string `json:"bid"`
	Branch string `json:"branch"`
	Cfg    branchmd.ProtectedBranchCfg
}

type ListProtectedBranchRespVO struct {
	ginutil.BaseResp
	Branches []ProtectedBranchVO
}
