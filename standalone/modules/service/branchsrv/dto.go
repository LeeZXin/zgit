package branchsrv

import (
	"fmt"
	"regexp"
	"zgit/pkg/apicode"
	"zgit/pkg/i18n"
	"zgit/standalone/modules/model/branchmd"
	"zgit/standalone/modules/model/usermd"
	"zgit/util"
)

var (
	validBranchPattern = regexp.MustCompile(`^\S{1,32}$`)
)

type InsertProtectedBranchReqDTO struct {
	RepoId   string
	Branch   string
	Cfg      branchmd.ProtectedBranchCfg
	Operator usermd.UserInfo
}

func (r *InsertProtectedBranchReqDTO) IsValid() error {
	if !validateRepoId(r.RepoId) {
		return util.InvalidArgsError()
	}
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !validateBranch(r.Branch) {
		fmt.Println(r.Branch)
		return util.InvalidArgsError()
	}
	if !validateOperator(r.Operator) {
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
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !validateBid(r.Bid) {
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
	Bid    string
	RepoId string
	Branch string
	Cfg    branchmd.ProtectedBranchCfg
}

func validateRepoId(repoId string) bool {
	return len(repoId) == 32
}

func validateBid(bid string) bool {
	return len(bid) == 32
}

func validateOperator(operator usermd.UserInfo) bool {
	return operator.Account != ""
}

func validateBranch(branch string) bool {
	return validBranchPattern.MatchString(branch)
}
