package pullrequestapi

type SubmitPullRequestReqVO struct {
	RepoId string `json:"repoId"`
	Target string `json:"target"`
	Head   string `json:"head"`
}

type ClosePullRequestReqVO struct {
	PrId string `json:"prId"`
}

type MergePullRequestReqVO struct {
	PrId string `json:"prId"`
}
