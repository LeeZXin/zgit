package projectmd

import (
	"context"
	"encoding/json"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"zgit/pkg/perm"
)

func GetByProjectId(ctx context.Context, projectId string) (Project, bool, error) {
	var ret Project
	b, err := xormutil.MustGetXormSession(ctx).Where("project_id = ?", projectId).Get(&ret)
	return ret, b, err
}

func InsertProject(ctx context.Context, reqDTO InsertProjectReqDTO) (Project, error) {
	ret := Project{
		ProjectId:   GenProjectId(),
		Name:        reqDTO.Name,
		ProjectDesc: reqDTO.Desc,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func GetProjectUserPermDetail(ctx context.Context, projectId, account string) (ProjectUserPermDetailDTO, bool, error) {
	pu, b, err := GetProjectUser(ctx, projectId, account)
	if err != nil || !b {
		return ProjectUserPermDetailDTO{}, b, err
	}
	if pu.GroupId == "" {
		return ProjectUserPermDetailDTO{}, true, nil
	}
	group, b, err := GetByGroupId(ctx, pu.GroupId)
	if err != nil || !b {
		return ProjectUserPermDetailDTO{}, b, err
	}
	ret, err := group.GetPermDetail()
	if err != nil {
		return ProjectUserPermDetailDTO{}, false, err
	}
	return ProjectUserPermDetailDTO{
		GroupId:    group.GroupId,
		IsAdmin:    group.IsAdmin,
		PermDetail: ret,
	}, true, nil
}

func InsertProjectUser(ctx context.Context, reqDTO InsertProjectUserReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ProjectUser{
		ProjectId: reqDTO.ProjectId,
		Account:   reqDTO.Account,
		GroupId:   reqDTO.GroupId,
	})
	return err
}

func UpdateProjectUser(ctx context.Context, reqDTO UpdateProjectUserReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("project_id = ?", reqDTO.ProjectId).
		And("account = ?", reqDTO.Account).
		Cols("group_id").
		Limit(1).
		Update(&ProjectUser{
			GroupId: reqDTO.GroupId,
		})
	return rows == 1, err
}

func DeleteProjectUser(ctx context.Context, projectId, account string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("project_id = ?", projectId).
		And("account = ?", account).
		Limit(1).
		Delete(new(ProjectUser))
	return rows == 1, err
}

func GetProjectUser(ctx context.Context, projectId, account string) (ProjectUser, bool, error) {
	ret := ProjectUser{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("project_id = ?", projectId).
		And("account = ?", account).
		Get(&ret)
	return ret, b, err
}

func InsertProjectUserGroup(ctx context.Context, reqDTO InsertProjectUserGroupReqDTO) (ProjectUserGroup, error) {
	m, _ := json.Marshal(reqDTO.PermDetail)
	ret := ProjectUserGroup{
		GroupId:   GenGroupId(),
		ProjectId: reqDTO.ProjectId,
		Name:      reqDTO.Name,
		Perm:      string(m),
		IsAdmin:   reqDTO.IsAdmin,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func GetByGroupId(ctx context.Context, groupId string) (ProjectUserGroup, bool, error) {
	ret := ProjectUserGroup{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("group_id = ?", groupId).
		Get(&ret)
	return ret, b, err
}

func UpdateProjectUserGroupName(ctx context.Context, groupId, name string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("group_id = ?", groupId).
		Cols("name").
		Limit(1).
		Update(&ProjectUserGroup{
			Name: name,
		})
	return rows == 1, err
}

func UpdateProjectUserGroupPerm(ctx context.Context, groupId string, detail perm.Detail) (bool, error) {
	m, _ := json.Marshal(detail)
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("group_id = ?", groupId).
		Cols("perm").
		Limit(1).
		Update(&ProjectUserGroup{
			Perm: string(m),
		})
	return rows == 1, err
}

func DeleteProjectUserGroup(ctx context.Context, groupId string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("group_id = ?", groupId).
		Limit(1).
		Delete(new(ProjectUserGroup))
	return rows == 1, err
}

func ExistProjectUser(ctx context.Context, projectId, groupId string) (bool, error) {
	// 走一下projectId索引
	return xormutil.MustGetXormSession(ctx).
		Where("project_id = ?", projectId).
		And("group_id = ?", groupId).
		Exist(new(ProjectUser))
}

func ListProjectUserGroup(ctx context.Context, projectId string) ([]ProjectUserGroup, error) {
	ret := make([]ProjectUserGroup, 0)
	err := xormutil.MustGetXormSession(ctx).Where("project_id = ?", projectId).Find(&ret)
	return ret, err
}
