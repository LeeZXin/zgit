package pullrequestsrv

import (
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
	if r.Operator.Account == "" {
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
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	if len(r.PrId) > 32 || len(r.PrId) == 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type MergePullRequestReqDTO struct {
	PrId     string
	Operator usermd.UserInfo
}

func (r *MergePullRequestReqDTO) IsValid() error {
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	if len(r.PrId) > 32 || len(r.PrId) == 0 {
		return util.InvalidArgsError()
	}
	return nil
}
