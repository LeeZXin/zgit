package repoapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"zgit/pkg/git"
)

type AllGitIgnoreTemplateListRespVO struct {
	ginutil.BaseResp
	Data []string `json:"data"`
}

type InitRepoReqVO struct {
	Name          string `json:"name"`
	Desc          string `json:"Desc"`
	RepoType      int    `json:"repoType"`
	CreateReadme  bool   `json:"createReadme"`
	ProjectId     string `json:"projectId"`
	GitIgnoreName string `json:"gitIgnoreName"`
	DefaultBranch string `json:"defaultBranch"`
}

type DeleteRepoReqVO struct {
	RepoId string `json:"repoId"`
}

type TreeRepoReqVO struct {
	RepoId  string `json:"repoId"`
	RefName string `json:"refName"`
	Dir     string `json:"dir"`
}

type CommitVO struct {
	Author        git.User
	Committer     git.User
	AuthoredDate  string
	CommittedDate string
	CommitMsg     string
	CommitId      string
	ShortId       string
}

type FileVO struct {
	Mode    string
	RawPath string
	Path    string
	Commit  CommitVO
}

type TreeVO struct {
	Files   []FileVO
	Limit   int
	HasMore bool
}

type TreeRepoRespVO struct {
	ginutil.BaseResp
	IsEmpty      bool
	ReadmeText   string
	RecentCommit CommitVO
	Tree         TreeVO
}
