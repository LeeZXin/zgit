package branchsrv

import (
	"github.com/LeeZXin/zsf-utils/bizerr"
	"zgit/pkg/apicode"
	"zgit/pkg/i18n"
	"zgit/standalone/modules/model/branchmd"
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
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if validateBranch(r.Branch) {
		return util.InvalidArgsError()
	}
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if len(r.Cfg.ReviewerList) > 50 {
		return util.InvalidArgsError()
	}
	if r.Cfg.ReviewCountWhenCreatePr < len(r.Cfg.ReviewerList) {
		return bizerr.NewBizErr(apicode.InvalidReviewCountWhenCreatePrCode.Int(), i18n.GetByKey(i18n.ProtectedBranchInvalidReviewCountWhenCreatePr))
	}
	return nil
}

type DeleteProtectedBranchReqDTO struct {
	RepoId   string
	Branch   string
	Operator usermd.UserInfo
}

func (r *DeleteProtectedBranchReqDTO) IsValid() error {
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if validateBranch(r.Branch) {
		return util.InvalidArgsError()
	}
	if !validateRepoId(r.RepoId) {
		return util.InvalidArgsError()
	}
	return nil
}

type ListProtectedBranchReqDTO struct {
	RepoId   string
	Operator usermd.UserInfo
}

func (r *ListProtectedBranchReqDTO) IsValid() error {
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !validateRepoId(r.RepoId) {
		return util.InvalidArgsError()
	}
	return nil
}

type ProtectedBranchDTO struct {
	RepoId string
	Branch string
	Cfg    branchmd.ProtectedBranchCfg
}

func validateRepoId(repoId string) bool {
	return len(repoId) == 32
}

func validateOperator(operator usermd.UserInfo) bool {
	return operator.Account != ""
}

func validateBranch(branch string) bool {
	return len(branch) <= 32 && len(branch) > 0
}
