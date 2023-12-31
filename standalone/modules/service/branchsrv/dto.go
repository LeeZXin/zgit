package branchsrv

import (
	"zgit/standalone/modules/model/branchmd"
	"zgit/standalone/modules/model/usermd"
	"zgit/util"
)

type InsertProtectedBranchReqDTO struct {
	RepoId   string
	Branch   string
	Operator usermd.UserInfo
}

func (r *InsertProtectedBranchReqDTO) IsValid() error {
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	if len(r.Branch) > 32 || len(r.Branch) == 0 {
		return util.InvalidArgsError()
	}
	if len(r.RepoId) > 32 || len(r.RepoId) == 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteProtectedBranchReqDTO struct {
	RepoId   string
	Branch   string
	Operator usermd.UserInfo
}

func (r *DeleteProtectedBranchReqDTO) IsValid() error {
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	if len(r.Branch) > 32 || len(r.Branch) == 0 {
		return util.InvalidArgsError()
	}
	if len(r.RepoId) > 32 || len(r.RepoId) == 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type ListProtectedBranchReqDTO struct {
	RepoId     string
	SearchName string
	Offset     int64
	Limit      int
	Operator   usermd.UserInfo
}

func (r *ListProtectedBranchReqDTO) IsValid() error {
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	if r.Offset < 0 {
		return util.InvalidArgsError()
	}
	if r.Limit < 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	if len(r.RepoId) > 32 || len(r.RepoId) == 0 {
		return util.InvalidArgsError()
	}
	if len(r.SearchName) > 255 {
		return util.InvalidArgsError()
	}
	return nil
}

type ListProtectedBranchRespDTO struct {
	Data       []branchmd.ProtectedBranch
	Cursor     int64
	TotalCount int64
}
