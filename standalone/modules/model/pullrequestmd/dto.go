package pullrequestmd

type InsertPullRequestReqDTO struct {
	RepoId   string
	Target   string
	Head     string
	CreateBy string
	PrStatus PrStatus
}
