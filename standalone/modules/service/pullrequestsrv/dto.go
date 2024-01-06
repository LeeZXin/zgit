package pullrequestsrv

import (
	"zgit/standalone/modules/model/pullrequestmd"
	"zgit/standalone/modules/model/repomd"
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
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !repomd.IsRepoIdValid(r.RepoId) {
		return util.InvalidArgsError()
	}
	if !util.ValidateRef(r.Target) {
		return util.InvalidArgsError()
	}
	if !util.ValidateRef(r.Head) {
		return util.InvalidArgsError()
	}
	return nil
}

type ClosePullRequestReqDTO struct {
	PrId     string
	Operator usermd.UserInfo
}

func (r *ClosePullRequestReqDTO) IsValid() error {
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !pullrequestmd.IsPrIdValid(r.PrId) {
		return util.InvalidArgsError()
	}
	return nil
}

type MergePullRequestReqDTO struct {
	PrId     string
	Operator usermd.UserInfo
}

func (r *MergePullRequestReqDTO) IsValid() error {
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !pullrequestmd.IsPrIdValid(r.PrId) {
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
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !pullrequestmd.IsPrIdValid(r.PrId) {
		return util.InvalidArgsError()
	}
	return nil
}
