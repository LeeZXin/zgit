package projectsrv

import (
	"zgit/pkg/perm"
	"zgit/standalone/modules/model/projectmd"
	"zgit/standalone/modules/model/usermd"
	"zgit/util"
)

type InsertProjectReqDTO struct {
	Name     string
	Desc     string
	Operator usermd.UserInfo
}

func (r *InsertProjectReqDTO) IsValid() error {
	if len(r.Name) == 0 || len(r.Name) > 64 {
		return util.InvalidArgsError()
	}
	if len(r.Desc) > 128 {
		return util.InvalidArgsError()
	}
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteProjectReqDTO struct {
	ProjectId string
	Operator  usermd.UserInfo
}

type DeleteProjectUserReqDTO struct {
	ProjectId string
	Account   string
	Operator  usermd.UserInfo
}

func (r *DeleteProjectUserReqDTO) IsValid() error {
	if !projectmd.IsProjectIdValid(r.ProjectId) {
		return util.InvalidArgsError()
	}
	if !usermd.IsUserAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type UpsertProjectUserReqDTO struct {
	ProjectId string
	Account   string
	GroupId   string
	Operator  usermd.UserInfo
}

func (r *UpsertProjectUserReqDTO) IsValid() error {
	if !projectmd.IsGroupIdValid(r.GroupId) {
		return util.InvalidArgsError()
	}
	if !projectmd.IsProjectIdValid(r.ProjectId) {
		return util.InvalidArgsError()
	}
	if !usermd.IsUserAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type InsertProjectUserGroupReqDTO struct {
	ProjectId string
	Name      string
	Perm      perm.Detail
	Operator  usermd.UserInfo
}

func (r *InsertProjectUserGroupReqDTO) IsValid() error {
	if !projectmd.IsGroupNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !projectmd.IsProjectIdValid(r.ProjectId) {
		return util.InvalidArgsError()
	}
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateProjectUserGroupNameReqDTO struct {
	GroupId  string
	Name     string
	Operator usermd.UserInfo
}

func (r *UpdateProjectUserGroupNameReqDTO) IsValid() error {
	if !projectmd.IsGroupNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !projectmd.IsGroupIdValid(r.GroupId) {
		return util.InvalidArgsError()
	}
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateProjectUserGroupPermReqDTO struct {
	GroupId  string
	Perm     perm.Detail
	Operator usermd.UserInfo
}

func (r *UpdateProjectUserGroupPermReqDTO) IsValid() error {
	if !projectmd.IsGroupIdValid(r.GroupId) {
		return util.InvalidArgsError()
	}
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteProjectUserGroupReqDTO struct {
	GroupId  string
	Operator usermd.UserInfo
}

func (r *DeleteProjectUserGroupReqDTO) IsValid() error {
	if !projectmd.IsGroupIdValid(r.GroupId) {
		return util.InvalidArgsError()
	}
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type ListProjectUserGroupReqDTO struct {
	ProjectId string
	Operator  usermd.UserInfo
}

func (r *ListProjectUserGroupReqDTO) IsValid() error {
	if !projectmd.IsProjectIdValid(r.ProjectId) {
		return util.InvalidArgsError()
	}
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type ProjectUserGroupDTO struct {
	GroupId   string
	ProjectId string
	Name      string
	Perm      perm.Detail
}
