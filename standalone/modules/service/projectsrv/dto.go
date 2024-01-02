package projectsrv

import (
	"zgit/pkg/perm"
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
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteProjectUserReqDTO struct {
	ProjectId string
	Account   string
	Operator  usermd.UserInfo
}

func (r *DeleteProjectUserReqDTO) IsValid() error {
	if len(r.ProjectId) == 0 || len(r.ProjectId) > 32 {
		return util.InvalidArgsError()
	}
	if len(r.Account) > 32 || len(r.Account) == 0 {
		return util.InvalidArgsError()
	}
	if r.Operator.Account == "" {
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
	if len(r.GroupId) > 32 {
		return util.InvalidArgsError()
	}
	if len(r.ProjectId) == 0 || len(r.ProjectId) > 32 {
		return util.InvalidArgsError()
	}
	if len(r.Account) > 32 || len(r.Account) == 0 {
		return util.InvalidArgsError()
	}
	if r.Operator.Account == "" {
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
	if len(r.Name) > 64 {
		return util.InvalidArgsError()
	}
	if len(r.ProjectId) == 0 || len(r.ProjectId) > 32 {
		return util.InvalidArgsError()
	}
	if r.Operator.Account == "" {
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
	if len(r.Name) == 0 || len(r.Name) > 64 {
		return util.InvalidArgsError()
	}
	if len(r.GroupId) == 0 || len(r.GroupId) > 32 {
		return util.InvalidArgsError()
	}
	if r.Operator.Account == "" {
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
	if len(r.GroupId) == 0 || len(r.GroupId) > 32 {
		return util.InvalidArgsError()
	}
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteProjectUserGroupReqDTO struct {
	GroupId  string
	Operator usermd.UserInfo
}

func (r *DeleteProjectUserGroupReqDTO) IsValid() error {
	if len(r.GroupId) == 0 || len(r.GroupId) > 32 {
		return util.InvalidArgsError()
	}
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	return nil
}

type ListProjectUserGroupReqDTO struct {
	ProjectId string
	Operator  usermd.UserInfo
}

func (r *ListProjectUserGroupReqDTO) IsValid() error {
	if len(r.ProjectId) == 0 || len(r.ProjectId) > 32 {
		return util.InvalidArgsError()
	}
	if r.Operator.Account == "" {
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
