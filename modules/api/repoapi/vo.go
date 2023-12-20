package repoapi

import "github.com/LeeZXin/zsf-utils/ginutil"

type AllGitIgnoreTemplateListRespVO struct {
	ginutil.BaseResp
	Data []string `json:"data"`
}

type InitRepoReqVO struct {
	RepoName      string `json:"repoName"`
	RepoDesc      string `json:"repoDesc"`
	RepoType      int    `json:"repoType"`
	CreateReadme  bool   `json:"createReadme"`
	GitIgnoreName string `json:"gitIgnoreName"`
	DefaultBranch string `json:"defaultBranch"`
}

type DeleteRepoReqVO struct {
	RepoId string `json:"repoId"`
}
