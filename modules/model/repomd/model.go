package repomd

import (
	"time"
	"zgit/pkg/i18n"
)

const (
	RepoTableName = "repo"
)

type RepoType int

const (
	InternalRepoType RepoType = iota
	PublicRepoType
)

var (
	repoTypeStringMap = map[RepoType]string{
		InternalRepoType: i18n.GetByKey(i18n.InternalRepoType),
		PublicRepoType:   i18n.GetByKey(i18n.PublicRepoType),
	}
)

func (t RepoType) String() string {
	ret, b := repoTypeStringMap[t]
	if b {
		return ret
	}
	return i18n.GetByKey(i18n.UnKnownRepoType)
}

func (t RepoType) IsValid() bool {
	_, b := repoTypeStringMap[t]
	return b
}

type RepoInfo struct {
	RepoId    string `json:"RepoId"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	UserId    string `json:"userId"`
	NodeId    string `json:"nodeId"`
	CorpId    string `json:"corpId"`
	ProjectId string `json:"projectId"`
	RepoType  int    `json:"repoType"`
	IsEmpty   bool   `json:"isEmpty"`
	TotalSize int64  `json:"totalSize"`
	GitSize   int64  `json:"gitSize"`
	LfsSize   int64  `json:"lfsSize"`
}

type Repo struct {
	Id            int64     `json:"id" xorm:"pk autoincr"`
	RepoId        string    `json:"repoId"`
	Name          string    `json:"name"`
	Path          string    `json:"path"`
	UserId        string    `json:"userId"`
	NodeId        string    `json:"nodeId"`
	CorpId        string    `json:"corpId"`
	ProjectId     string    `json:"projectId"`
	RepoDesc      string    `json:"repoDesc"`
	DefaultBranch string    `json:"defaultBranch"`
	RepoType      int       `json:"repoType"`
	IsEmpty       bool      `json:"isEmpty"`
	TotalSize     int64     `json:"totalSize"`
	GitSize       int64     `json:"gitSize"`
	LfsSize       int64     `json:"lfsSize"`
	Created       time.Time `json:"created" xorm:"created"`
	Updated       time.Time `json:"updated" xorm:"updated"`
}

func (*Repo) TableName() string {
	return RepoTableName
}

func (r *Repo) ToRepoInfo() RepoInfo {
	return RepoInfo{
		RepoId:    r.RepoId,
		Name:      r.Name,
		Path:      r.Path,
		UserId:    r.UserId,
		NodeId:    r.NodeId,
		CorpId:    r.CorpId,
		ProjectId: r.ProjectId,
		RepoType:  r.RepoType,
		IsEmpty:   r.IsEmpty,
		TotalSize: r.TotalSize,
		GitSize:   r.GitSize,
		LfsSize:   r.LfsSize,
	}
}
