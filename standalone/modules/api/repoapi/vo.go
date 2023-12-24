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
	RepoPath string `json:"repoPath"`
}

type TreeRepoReqVO struct {
	RepoPath string `json:"repoPath"`
	RefName  string `json:"refName"`
	Dir      string `json:"dir"`
}

type EntriesRepoReqVO struct {
	RepoPath string `json:"repoPath"`
	RefName  string `json:"refName"`
	Dir      string `json:"dir"`
	Offset   int    `json:"offset"`
}

type ListRepoReqVO struct {
	Offset     int64  `json:"offset"`
	Limit      int    `json:"limit"`
	SearchName string `json:"searchName"`
	ProjectId  string `json:"projectId"`
}

type ListRepoRespVO struct {
	ginutil.BaseResp
	RepoList   []RepoVO `json:"repoList"`
	TotalCount int64    `json:"totalCount"`
	Cursor     int64    `json:"cursor"`
	Limit      int      `json:"limit"`
}

type CommitVO struct {
	Author        git.User `json:"author"`
	Committer     git.User `json:"committer"`
	AuthoredDate  string   `json:"authoredDate"`
	CommittedDate string   `json:"committedDate"`
	CommitMsg     string   `json:"commitMsg"`
	CommitId      string   `json:"commitId"`
	ShortId       string   `json:"shortId"`
}

type FileVO struct {
	Mode    string   `json:"mode"`
	RawPath string   `json:"rawPath"`
	Path    string   `json:"path"`
	Commit  CommitVO `json:"commit"`
}

type TreeVO struct {
	Files   []FileVO `json:"files"`
	Limit   int      `json:"limit"`
	Offset  int      `json:"offset"`
	HasMore bool     `json:"hasMore"`
}

type TreeRepoRespVO struct {
	ginutil.BaseResp
	IsEmpty      bool     `json:"isEmpty"`
	ReadmeText   string   `json:"readmeText"`
	RecentCommit CommitVO `json:"recentCommit"`
	Tree         TreeVO   `json:"tree"`
}

type RepoVO struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Author    string `json:"author"`
	ProjectId string `json:"projectId"`
	RepoType  string `json:"repoType"`
	IsEmpty   bool   `json:"isEmpty"`
	TotalSize int64  `json:"totalSize"`
	WikiSize  int64  `json:"wikiSize"`
	GitSize   int64  `json:"gitSize"`
	LfsSize   int64  `json:"lfsSize"`
	Created   string `json:"created"`
}

type CatFileReqVO struct {
	RepoPath string `json:"repoPath"`
	RefName  string `json:"refName"`
	Dir      string `json:"dir"`
	FileName string `json:"fileName"`
}

type CatFileRespVO struct {
	ginutil.BaseResp
	Mode    string `json:"mode"`
	Content string `json:"content"`
}
