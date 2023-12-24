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
	PrivateRepoType
)

var (
	repoTypeStringMap = map[RepoType]string{
		InternalRepoType: i18n.GetByKey(i18n.InternalRepoType),
		PublicRepoType:   i18n.GetByKey(i18n.PublicRepoType),
		PrivateRepoType:  i18n.GetByKey(i18n.PrivateRepoType),
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
	Name      string `json:"name"`
	Path      string `json:"path"`
	Author    string `json:"author"`
	ProjectId string `json:"projectId"`
	RepoType  int    `json:"repoType"`
	IsEmpty   bool   `json:"isEmpty"`
	TotalSize int64  `json:"totalSize"`
	WikiSize  int64  `json:"wikiSize"`
	GitSize   int64  `json:"gitSize"`
	LfsSize   int64  `json:"lfsSize"`
}

type Repo struct {
	Id            int64     `json:"id" xorm:"pk autoincr"`
	Path          string    `json:"path"`
	Name          string    `json:"name"`
	Author        string    `json:"author"`
	ProjectId     string    `json:"projectId"`
	RepoDesc      string    `json:"repoDesc"`
	DefaultBranch string    `json:"defaultBranch"`
	RepoType      int       `json:"repoType"`
	IsEmpty       bool      `json:"isEmpty"`
	TotalSize     int64     `json:"totalSize"`
	WikiSize      int64     `json:"wikiSize"`
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
		Name:      r.Name,
		Path:      r.Path,
		Author:    r.Author,
		ProjectId: r.ProjectId,
		RepoType:  r.RepoType,
		IsEmpty:   r.IsEmpty,
		TotalSize: r.TotalSize,
		GitSize:   r.GitSize,
		LfsSize:   r.LfsSize,
		WikiSize:  r.WikiSize,
	}
}
