package branchsrv

import (
	"zgit/pkg/apicode"
	"zgit/pkg/i18n"
	"zgit/standalone/modules/model/branchmd"
	"zgit/standalone/modules/model/repomd"
	"zgit/standalone/modules/model/usermd"
	"zgit/util"
)

type InsertProtectedBranchReqDTO struct {
	RepoId   string
	Branch   string
	Cfg      branchmd.ProtectedBranchCfg
	Operator usermd.UserInfo
}

func (r *InsertProtectedBranchReqDTO) IsValid() error {
	if !repomd.IsRepoIdValid(r.RepoId) {
		return util.InvalidArgsError()
	}
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !branchmd.IsWildcardBranchValid(r.Branch) {
		return util.InvalidArgsError()
	}
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if len(r.Cfg.ReviewerList) > 50 {
		return util.InvalidArgsError()
	}
	if r.Cfg.ReviewCountWhenCreatePr < len(r.Cfg.ReviewerList) {
		return util.NewBizErr(apicode.InvalidReviewCountWhenCreatePrCode, i18n.ProtectedBranchInvalidReviewCountWhenCreatePr)
	}
	return nil
}

type DeleteProtectedBranchReqDTO struct {
	Bid      string
	Operator usermd.UserInfo
}

func (r *DeleteProtectedBranchReqDTO) IsValid() error {
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !branchmd.IsBidValid(r.Bid) {
		return util.InvalidArgsError()
	}
	return nil
}

type ListProtectedBranchReqDTO struct {
	RepoId   string
	Operator usermd.UserInfo
}

func (r *ListProtectedBranchReqDTO) IsValid() error {
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !repomd.IsRepoIdValid(r.RepoId) {
		return util.InvalidArgsError()
	}
	return nil
}

type ProtectedBranchDTO struct {
	Bid    string
	RepoId string
	Branch string
	Cfg    branchmd.ProtectedBranchCfg
}
