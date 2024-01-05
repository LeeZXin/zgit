package branchmd

type InsertProtectedBranchReqDTO struct {
	RepoId string
	Branch string
	Cfg    ProtectedBranchCfg
}

type ProtectedBranchDTO struct {
	Bid    string
	RepoId string
	Branch string
	Cfg    ProtectedBranchCfg
}
