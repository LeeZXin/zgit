package pullrequestapi

import "zgit/standalone/modules/model/pullrequestmd"

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

type ReviewPullRequestReqVO struct {
	PrId      string                     `json:"prId"`
	Status    pullrequestmd.ReviewStatus `json:"status"`
	ReviewMsg string                     `json:"reviewMsg"`
}
