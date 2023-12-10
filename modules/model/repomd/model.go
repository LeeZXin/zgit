package repomd

type RepoInfo struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	OwnerId   string `json:"ownerId"`
	ClusterId string `json:"clusterId"`
}
