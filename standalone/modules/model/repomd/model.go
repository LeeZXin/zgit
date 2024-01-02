package repomd

import (
	"time"
)

const (
	RepoTableName = "repo"
)

type Repo struct {
	Id            int64     `json:"id" xorm:"pk autoincr"`
	RepoId        string    `json:"repoId"`
	Path          string    `json:"path"`
	Name          string    `json:"name"`
	Author        string    `json:"author"`
	ProjectId     string    `json:"projectId"`
	RepoDesc      string    `json:"repoDesc"`
	DefaultBranch string    `json:"defaultBranch"`
	RepoType      int       `json:"repoType"`
	RepoStatus    int       `json:"repoStatus"`
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
		RepoId:    r.RepoId,
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
