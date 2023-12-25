package branchmd

type ListProtectedBranchReqDTO struct {
	RepoPath   string
	SearchName string
	Offset     int64
	Limit      int
}
