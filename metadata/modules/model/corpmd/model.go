package corpmd

import (
	"time"
)

const (
	CorpTableName = "corp"
)

type CorpInfo struct {
	CorpId     string `json:"corpId"`
	Name       string `json:"name"`
	NodeId     string `json:"nodeId"`
	RepoCount  int    `json:"repoCount"`
	RepoLimit  int    `json:"repoLimit"`
	MaxLfsSize int    `json:"maxLfsSize"`
	MaxGitSize int    `json:"maxGitSize"`
}

type Corp struct {
	Id         int64     `json:"id" xorm:"pk autoincr"`
	CorpId     string    `json:"corpId"`
	Name       string    `json:"name"`
	NodeId     string    `json:"nodeId"`
	CorpDesc   string    `json:"corpDesc"`
	RepoCount  int       `json:"repoCount"`
	RepoLimit  int       `json:"repoLimit"`
	MaxGitSize int       `json:"maxGitSize"`
	MaxLfsSize int       `json:"maxLfsSize"`
	Created    time.Time `json:"created" xorm:"created"`
	Updated    time.Time `json:"updated" xorm:"updated"`
}

func (*Corp) TableName() string {
	return CorpTableName
}

func (c *Corp) ToCorpInfo() CorpInfo {
	return CorpInfo{
		CorpId:     c.CorpId,
		Name:       c.Name,
		NodeId:     c.NodeId,
		RepoCount:  c.RepoCount,
		RepoLimit:  c.RepoLimit,
		MaxLfsSize: c.MaxLfsSize,
		MaxGitSize: c.MaxGitSize,
	}
}
