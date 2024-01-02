package projectmd

import "zgit/pkg/perm"

type InsertProjectReqDTO struct {
	Name string
	Desc string
}

type InsertProjectUserReqDTO struct {
	ProjectId string
	Account   string
	GroupId   string
}

type UpdateProjectUserReqDTO struct {
	ProjectId string
	Account   string
	GroupId   string
}

type InsertProjectUserGroupReqDTO struct {
	Name       string
	ProjectId  string
	PermDetail perm.Detail
	IsAdmin    bool
}

type ProjectUserPermDetailDTO struct {
	GroupId    string
	IsAdmin    bool
	PermDetail perm.Detail
}
