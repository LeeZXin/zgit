package projectmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func GetByProjectId(ctx context.Context, projectId string) (Project, bool, error) {
	var ret Project
	b, err := xormutil.MustGetXormSession(ctx).Where("project_id = ?", projectId).Get(&ret)
	return ret, b, err
}

func InsertProject(ctx context.Context, reqDTO InsertProjectReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&Project{
		ProjectId:   GenProjectId(),
		Name:        reqDTO.Name,
		ProjectDesc: reqDTO.Desc,
	})
	return err
}

func ProjectUserExists(ctx context.Context, projectId, account string) (bool, error) {
	return xormutil.MustGetXormSession(ctx).Where("project_id = ?", projectId).
		And("account = ?", account).
		Exist(new(ProjectUser))
}

func InsertProjectUser(ctx context.Context, projectId, account string) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ProjectUser{
		ProjectId: projectId,
		Account:   account,
	})
	return err
}

func DeleteProjectUser(ctx context.Context, projectId, account string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("project_id = ?", projectId).
		And("account = ?", account).
		Limit(1).
		Delete(new(ProjectUser))
	return rows == 1, err
}

func ListProjectUserByAccount(ctx context.Context, account string) ([]ProjectUser, error) {
	ret := make([]ProjectUser, 0)
	err := xormutil.MustGetXormSession(ctx).Where("account = ?", account).Find(&ret)
	return ret, err
}
