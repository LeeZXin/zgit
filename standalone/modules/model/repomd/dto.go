package repomd

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
