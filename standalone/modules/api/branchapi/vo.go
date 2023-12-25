package branchapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
)

type InsertProtectedBranchReqVO struct {
	RepoPath string `json:"repoPath"`
	Branch   string `json:"branch"`
}

type DeleteProtectedBranchReqVO struct {
	RepoPath string `json:"repoPath"`
	Branch   string `json:"branch"`
}

type ListProtectedBranchReqVO struct {
	RepoPath   string `json:"repoPath"`
	SearchName string `json:"searchName"`
	Offset     int64  `json:"offset"`
	Limit      int    `json:"limit"`
}

type ProtectedBranchVO struct {
	Branch  string `json:"branch"`
	Created string `json:"created"`
}

type ListProtectedBranchRespVO struct {
	ginutil.BaseResp
	Branches   []ProtectedBranchVO
	Cursor     int64
	TotalCount int64
}
