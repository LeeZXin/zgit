package pullrequestsrv

import (
	"zgit/standalone/modules/model/pullrequestmd"
	"zgit/standalone/modules/model/usermd"
	"zgit/util"
)

type SubmitPullRequestReqDTO struct {
	RepoId   string
	Target   string
	Head     string
	Operator usermd.UserInfo
}

func (r *SubmitPullRequestReqDTO) IsValid() error {
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if len(r.RepoId) > 32 || len(r.RepoId) == 0 {
		return util.InvalidArgsError()
	}
	if len(r.Target) > 128 || len(r.Target) == 0 {
		return util.InvalidArgsError()
	}
	if len(r.Head) > 128 || len(r.Head) == 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type ClosePullRequestReqDTO struct {
	PrId     string
	Operator usermd.UserInfo
}

func (r *ClosePullRequestReqDTO) IsValid() error {
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !validatePrId(r.PrId) {
		return util.InvalidArgsError()
	}
	return nil
}

type MergePullRequestReqDTO struct {
	PrId     string
	Operator usermd.UserInfo
}

func (r *MergePullRequestReqDTO) IsValid() error {
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !validatePrId(r.PrId) {
		return util.InvalidArgsError()
	}
	return nil
}

type ReviewPullRequestReqDTO struct {
	PrId      string
	Status    pullrequestmd.ReviewStatus
	ReviewMsg string
	Operator  usermd.UserInfo
}

func (r *ReviewPullRequestReqDTO) IsValid() error {
	if len(r.ReviewMsg) > 255 {
		return util.InvalidArgsError()
	}
	if !r.Status.IsValid() {
		return util.InvalidArgsError()
	}
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !validatePrId(r.PrId) {
		return util.InvalidArgsError()
	}
	return nil
}

func validateOperator(operator usermd.UserInfo) bool {
	return operator.Account != ""
}

func validatePrId(prId string) bool {
	return len(prId) == 32
}
