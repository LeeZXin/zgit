package branchmd

type ListProtectedBranchReqDTO struct {
	RepoId     string
	SearchName string
	Offset     int64
	Limit      int
}
