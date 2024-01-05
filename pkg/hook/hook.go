package hook

const (
	ApiPreReceiveUrl  = "api/internal/hook/pre-receive"
	ApiPostReceiveUrl = "api/internal/hook/post-receive"
)

type RevInfo struct {
	OldCommitId string `json:"oldCommitId"`
	NewCommitId string `json:"newCommitId"`
	RefName     string `json:"refName"`
}

type Opts struct {
	RevInfoList                  []RevInfo `json:"revInfoList"`
	RepoId                       string    `json:"repoId"`
	PrId                         string    `json:"prId"`
	PusherId                     string    `json:"pusherId"`
	ObjectDirectory              string    `json:"objectDirectory"`
	AlternativeObjectDirectories string    `json:"alternativeObjectDirectories"`
	QuarantinePath               string    `json:"quarantinePath"`
}
