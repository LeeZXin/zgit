package repomd

type RepoInfo struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	RelativePath string `json:"relativePath"`
	OwnerId      string `json:"ownerId"`
	ClusterId    string `json:"clusterId"`
}
