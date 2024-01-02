package repomd

import "zgit/pkg/i18n"

type InsertRepoReqDTO struct {
	Name          string   `json:"name"`
	Path          string   `json:"path"`
	Author        string   `json:"author"`
	ProjectId     string   `json:"projectId"`
	RepoDesc      string   `json:"repoDesc"`
	DefaultBranch string   `json:"defaultBranch"`
	RepoType      RepoType `json:"repoType"`
	IsEmpty       bool     `json:"isEmpty"`
	TotalSize     int64    `json:"totalSize"`
	GitSize       int64    `json:"gitSize"`
	LfsSize       int64    `json:"lfsSize"`
}

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

func (t RepoType) Int() int {
	return int(t)
}

func (t RepoType) Readable() string {
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
	RepoId    string `json:"repoId"`
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

type RepoStatus int

const (
	OpenRepoStatus RepoStatus = iota
	ClosedRepoStatus
	DeletedRepoStatus
)

func (s RepoStatus) Int() int {
	return int(s)
}

func (s RepoStatus) Readable() string {
	switch s {
	case OpenRepoStatus:
		return i18n.GetByKey(i18n.RepoOpenStatus)
	case ClosedRepoStatus:
		return i18n.GetByKey(i18n.RepoClosedStatus)
	case DeletedRepoStatus:
		return i18n.GetByKey(i18n.RepoDeletedStatus)
	default:
		return i18n.GetByKey(i18n.RepoUnknownStatus)
	}
}
