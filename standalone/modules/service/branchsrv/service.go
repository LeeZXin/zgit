package branchsrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"zgit/pkg/apicode"
	"zgit/pkg/i18n"
	"zgit/standalone/modules/model/branchmd"
	"zgit/standalone/modules/model/projectmd"
	"zgit/standalone/modules/model/repomd"
	"zgit/standalone/modules/model/usermd"
	"zgit/util"
)

func InsertProtectedBranch(ctx context.Context, reqDTO InsertProtectedBranchReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	err := checkPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, b, err := branchmd.GetProtectedBranch(ctx, reqDTO.RepoId, reqDTO.Branch)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if b {
		return util.AlreadyExistsError()
	}
	for _, account := range reqDTO.Cfg.ReviewerList {
		_, b, err = usermd.GetByAccount(ctx, account)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError()
		}
		if !b {
			return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), fmt.Sprintf(i18n.GetByKey(i18n.UserAccountNotFoundWarnFormat), account))
		}
	}
	if err = branchmd.InsertProtectedBranch(ctx, branchmd.InsertProtectedBranchReqDTO{
		RepoId: reqDTO.RepoId,
		Branch: reqDTO.Branch,
		Cfg:    reqDTO.Cfg,
	}); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	return nil
}

func DeleteProtectedBranch(ctx context.Context, reqDTO DeleteProtectedBranchReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	err := checkPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return err
	}
	branch, b, err := branchmd.GetProtectedBranch(ctx, reqDTO.RepoId, reqDTO.Branch)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if !b {
		return util.InvalidArgsError()
	}
	_, err = branchmd.DeleteProtectedBranch(ctx, branch)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	return nil
}

func ListProtectedBranch(ctx context.Context, reqDTO ListProtectedBranchReqDTO) ([]ProtectedBranchDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	err := checkPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	branchList, err := branchmd.ListProtectedBranch(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError()
	}
	ret, _ := listutil.Map(branchList, func(t branchmd.ProtectedBranchDTO) (ProtectedBranchDTO, error) {
		return ProtectedBranchDTO{
			RepoId: t.RepoId,
			Branch: t.Branch,
			Cfg:    t.Cfg,
		}, nil
	})
	return ret, nil
}

// checkPerm 检查权限 只需检查是否是项目管理员
func checkPerm(ctx context.Context, repoId string, operator usermd.UserInfo) error {
	// 检查仓库是否存在
	repo, b, err := repomd.GetByRepoId(ctx, repoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 如果是系统管理员有所有权限
	if operator.IsAdmin {
		return nil
	}
	// 如果不是 检查用户组权限
	permDetail, b, err := projectmd.GetProjectUserPermDetail(ctx, repo.ProjectId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if !b {
		return util.UnauthorizedError()
	}
	if !permDetail.PermDetail.GetRepoPerm(repoId).CanHandleProtectedBranch {
		return util.UnauthorizedError()
	}
	return nil
}
