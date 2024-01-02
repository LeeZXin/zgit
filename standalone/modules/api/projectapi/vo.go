package projectapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"zgit/pkg/perm"
)

type InsertProjectReqVO struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type UpsertProjectUserReqVO struct {
	ProjectId string `json:"projectId"`
	Account   string `json:"account"`
	GroupId   string `json:"groupId"`
}

type DeleteProjectUserReqVO struct {
	ProjectId string `json:"projectId"`
	Account   string `json:"account"`
}

type InsertProjectUserGroupReqVO struct {
	ProjectId string      `json:"projectId"`
	Name      string      `json:"name"`
	Perm      perm.Detail `json:"perm"`
}

type UpdateProjectUserGroupNameReqVO struct {
	GroupId string `json:"groupId"`
	Name    string `json:"name"`
}

type UpdateProjectUserGroupPermReqVO struct {
	GroupId string      `json:"groupId"`
	Perm    perm.Detail `json:"perm"`
}

type DeleteProjectUserGroupReqVO struct {
	GroupId string `json:"groupId"`
}

type ListProjectUserGroupReqVO struct {
	ProjectId string `json:"projectId"`
}

type ProjectUserGroupVO struct {
	GroupId   string      `json:"groupId"`
	ProjectId string      `json:"projectId"`
	Name      string      `json:"name"`
	Perm      perm.Detail `json:"perm"`
}

type ListProjectUserGroupRespVO struct {
	ginutil.BaseResp
	Data []ProjectUserGroupVO
}
