package corpmd

type InsertCorpReqDTO struct {
	CorpId    string `json:"corpId"`
	Name      string `json:"name"`
	CorpDesc  string `json:"corpDesc"`
	RepoLimit int    `json:"repoLimit"`
}
