package projectsrv

import (
	"context"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"zgit/pkg/apicode"
	"zgit/pkg/i18n"
	"zgit/pkg/perm"
	"zgit/standalone/modules/model/projectmd"
	"zgit/standalone/modules/model/usermd"
	"zgit/standalone/modules/service/cfgsrv"
	"zgit/util"
)

// InsertProject 创建项目
func InsertProject(ctx context.Context, reqDTO InsertProjectReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 判断是否是系统管理员
	if !reqDTO.Operator.IsAdmin {
		// 判断是否允许用户自行创建项目
		sysCfg, err := cfgsrv.GetSysCfgWithCache(ctx)
		if err != nil {
			return err
		}
		if !sysCfg.AllowUserCreateProject {
			return util.UnauthorizedError()
		}
	}
	err := mysqlstore.WithTx(ctx, func(ctx context.Context) error {
		// 创建项目
		pu, err := projectmd.InsertProject(ctx, projectmd.InsertProjectReqDTO{
			Name: reqDTO.Name,
			Desc: reqDTO.Desc,
		})
		if err != nil {
			return err
		}
		// 创建管理员组
		group, err := projectmd.InsertProjectUserGroup(ctx, projectmd.InsertProjectUserGroupReqDTO{
			Name:       i18n.GetByKey(i18n.ProjectAdminUserGroupName),
			ProjectId:  pu.ProjectId,
			PermDetail: perm.DefaultPermDetail,
			IsAdmin:    true,
		})
		if err != nil {
			return err
		}
		// 创建关联关系 CreatorProjectUserType
		err = projectmd.InsertProjectUser(ctx, projectmd.InsertProjectUserReqDTO{
			ProjectId: pu.ProjectId,
			Account:   reqDTO.Operator.Account,
			GroupId:   group.GroupId,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	return nil
}

func DeleteProjectUser(ctx context.Context, reqDTO DeleteProjectUserReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	if err := checkProjectUserPerm(ctx, reqDTO.ProjectId, reqDTO.Operator); err != nil {
		return err
	}
	_, b, err := projectmd.GetProjectUser(ctx, reqDTO.ProjectId, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if !b {
		return util.InvalidArgsError()
	}
	_, err = projectmd.DeleteProjectUser(ctx, reqDTO.ProjectId, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	return nil
}

func UpsertProjectUser(ctx context.Context, reqDTO UpsertProjectUserReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	if err := checkProjectUserPerm(ctx, reqDTO.ProjectId, reqDTO.Operator); err != nil {
		return err
	}
	// 校验groupId是否正确
	if reqDTO.GroupId != "" {
		group, b, err := projectmd.GetByGroupId(ctx, reqDTO.GroupId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError()
		}
		// projectId不匹配
		if !b || group.ProjectId != reqDTO.ProjectId {
			return util.InvalidArgsError()
		}
	}
	_, b, err := projectmd.GetProjectUser(ctx, reqDTO.ProjectId, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if !b {
		// 校验账号是否存在
		_, b, err = usermd.GetByAccount(ctx, reqDTO.Account)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError()
		}
		if !b {
			return util.InvalidArgsError()
		}
		// 不存在则插入
		err = projectmd.InsertProjectUser(ctx, projectmd.InsertProjectUserReqDTO{
			ProjectId: reqDTO.ProjectId,
			Account:   reqDTO.Account,
			GroupId:   reqDTO.GroupId,
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError()
		}
	} else {
		_, err = projectmd.UpdateProjectUser(ctx, projectmd.UpdateProjectUserReqDTO{
			ProjectId: reqDTO.ProjectId,
			Account:   reqDTO.Account,
			GroupId:   reqDTO.GroupId,
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError()
		}
	}
	return nil
}

func checkProjectUserPerm(ctx context.Context, projectId string, operator usermd.UserInfo) error {
	// 判断权限
	if !operator.IsAdmin {
		pu, b, err := projectmd.GetProjectUser(ctx, projectId, operator.Account)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError()
		}
		// 不存在或不是管理员角色
		if !b || pu.GroupId == "" {
			return util.UnauthorizedError()
		}
		group, b, err := projectmd.GetByGroupId(ctx, pu.GroupId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError()
		}
		// groupId为空 用户组应该存在 不存在就是bug
		if !b {
			logger.Logger.WithContext(ctx).Errorf("check project: %s, user: %s found not group: %s", projectId, operator.Account, pu.GroupId)
			return util.InternalError()
		}
		// 不是项目管理员组
		if !group.IsAdmin {
			return util.UnauthorizedError()
		}
	} else {
		// 校验projectId是否存在
		_, b, err := projectmd.GetByProjectId(ctx, projectId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError()
		}
		if !b {
			return util.InvalidArgsError()
		}
	}
	return nil
}

func checkProjectUserPermByGroupId(ctx context.Context, operator usermd.UserInfo, groupId string) (projectmd.ProjectUserGroup, error) {
	// 检查权限
	group, b, err := projectmd.GetByGroupId(ctx, groupId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return projectmd.ProjectUserGroup{}, util.InternalError()
	}
	if !b {
		return projectmd.ProjectUserGroup{}, util.InvalidArgsError()
	}
	// 系统管理员有权限
	if operator.IsAdmin {
		return group, nil
	}
	// 非系统管理员检查项目管理员权限
	pu, b, err := projectmd.GetProjectUser(ctx, group.ProjectId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return group, util.InternalError()
	}
	// 不存在或不是管理员角色
	if !b || pu.GroupId == "" {
		return group, util.UnauthorizedError()
	}
	if !group.IsAdmin {
		return group, util.UnauthorizedError()
	}
	return group, nil
}

func InsertProjectUserGroup(ctx context.Context, reqDTO InsertProjectUserGroupReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	if err := checkProjectUserPerm(ctx, reqDTO.ProjectId, reqDTO.Operator); err != nil {
		return err
	}
	if _, err := projectmd.InsertProjectUserGroup(ctx, projectmd.InsertProjectUserGroupReqDTO{
		Name:       reqDTO.Name,
		ProjectId:  reqDTO.ProjectId,
		PermDetail: reqDTO.Perm,
	}); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	return nil
}

func UpdateProjectUserGroupName(ctx context.Context, reqDTO UpdateProjectUserGroupNameReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	group, err := checkProjectUserPermByGroupId(ctx, reqDTO.Operator, reqDTO.GroupId)
	if err != nil {
		return err
	}
	// 管理员项目组无法编辑权限
	if group.IsAdmin {
		return bizerr.NewBizErr(apicode.CannotUpdateProjectUserAdminGroupCode.Int(), i18n.GetByKey(i18n.ProjectUserGroupUpdateAdminNotAllow))
	}
	if _, err := projectmd.UpdateProjectUserGroupName(ctx, reqDTO.GroupId, reqDTO.Name); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	return nil
}

func UpdateProjectUserGroupPerm(ctx context.Context, reqDTO UpdateProjectUserGroupPermReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	group, err := checkProjectUserPermByGroupId(ctx, reqDTO.Operator, reqDTO.GroupId)
	if err != nil {
		return err
	}
	// 管理员项目组无法编辑权限
	if group.IsAdmin {
		return bizerr.NewBizErr(apicode.CannotUpdateProjectUserAdminGroupCode.Int(), i18n.GetByKey(i18n.ProjectUserGroupUpdateAdminNotAllow))
	}
	if _, err = projectmd.UpdateProjectUserGroupPerm(ctx, reqDTO.GroupId, reqDTO.Perm); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	return nil
}

func DeleteProjectUserGroup(ctx context.Context, reqDTO DeleteProjectUserGroupReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	group, err := checkProjectUserPermByGroupId(ctx, reqDTO.Operator, reqDTO.GroupId)
	if err != nil {
		return err
	}
	b, err := projectmd.ExistProjectUser(ctx, group.ProjectId, reqDTO.GroupId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	// 存在属于该groupId的用户
	if b {
		return bizerr.NewBizErr(apicode.ProjectUserGroupHasUserWhenDelCode.Int(), i18n.GetByKey(i18n.ProjectUserGroupHasUserWhenDel))
	}
	if _, err = projectmd.DeleteProjectUserGroup(ctx, reqDTO.GroupId); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	return nil
}

func ListProjectUserGroup(ctx context.Context, reqDTO ListProjectUserGroupReqDTO) ([]ProjectUserGroupDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	if err := checkProjectUserPerm(ctx, reqDTO.ProjectId, reqDTO.Operator); err != nil {
		return nil, err
	}
	groups, err := projectmd.ListProjectUserGroup(ctx, reqDTO.ProjectId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError()
	}
	return listutil.Map(groups, func(t projectmd.ProjectUserGroup) (ProjectUserGroupDTO, error) {
		detail, _ := t.GetPermDetail()
		return ProjectUserGroupDTO{
			GroupId:   t.GroupId,
			ProjectId: t.ProjectId,
			Name:      t.Name,
			Perm:      detail,
		}, nil
	})
}
