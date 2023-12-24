package corpsrv

type CorpInfoDTO struct {
	CorpId     string `json:"corpId"`
	Name       string `json:"name"`
	NodeId     string `json:"nodeId"`
	RepoCount  int    `json:"repoCount"`
	RepoLimit  int    `json:"repoLimit"`
	MaxLfsSize int    `json:"maxLfsSize"`
	MaxGitSize int    `json:"maxGitSize"`
}
