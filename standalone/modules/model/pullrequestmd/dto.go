package pullrequestmd

type InsertPullRequestReqDTO struct {
	RepoId   string
	Target   string
	Head     string
	CreateBy string
	PrStatus PrStatus
}

type InsertReviewReqDTO struct {
	PrId      string
	ReviewMsg string
	Status    ReviewStatus
	Reviewer  string
}

type UpdateReviewReqDTO struct {
	Rid    string
	Status ReviewStatus
}
