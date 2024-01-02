package branchsrv

import (
	"context"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
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
	_, err = branchmd.InsertProtectedBranch(ctx, reqDTO.RepoId, reqDTO.Branch)
	if err != nil {
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

func ListProtectedBranch(ctx context.Context, reqDTO ListProtectedBranchReqDTO) (ListProtectedBranchRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return ListProtectedBranchRespDTO{}, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	err := checkPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return ListProtectedBranchRespDTO{}, err
	}
	daoDto := branchmd.ListProtectedBranchReqDTO{
		RepoId:     reqDTO.RepoId,
		SearchName: reqDTO.SearchName,
		Offset:     reqDTO.Offset,
		Limit:      reqDTO.Limit,
	}
	ret := ListProtectedBranchRespDTO{}
	ret.TotalCount, err = branchmd.CountProtectedBranch(ctx, daoDto)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return ListProtectedBranchRespDTO{}, util.InternalError()
	}
	if ret.TotalCount > 0 {
		ret.Data, err = branchmd.ListProtectedBranch(ctx, daoDto)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return ListProtectedBranchRespDTO{}, util.InternalError()
		}
		ret.Cursor = ret.Data[len(ret.Data)-1].Id
	}
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
