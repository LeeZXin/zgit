package branchmd

import (
	"encoding/json"
	"time"
)

const (
	ProtectedBranchTableName = "protected_branch"
)

type ProtectedBranch struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	Bid     string    `json:"bid"`
	Branch  string    `json:"branch"`
	RepoId  string    `json:"repoId"`
	Cfg     string    `json:"cfg"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*ProtectedBranch) TableName() string {
	return ProtectedBranchTableName
}

type ProtectedBranchCfg struct {
	// 合并请求时代码评审数量大于该数量才能合并
	ReviewCountWhenCreatePr int `json:"ReviewCountWhenCreatePr"`
	// 代码评审员名单
	ReviewerList []string `json:"reviewerList"`
	// 可直接推送名单
	DirectPushList []string `json:"directPushList"`
}

func (c *ProtectedBranchCfg) ToString() string {
	m, _ := json.Marshal(c)
	return string(m)
}
