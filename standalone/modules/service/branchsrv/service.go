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
	err := checkAuth(ctx, reqDTO.RepoPath, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, b, err := branchmd.GetProtectedBranch(ctx, reqDTO.RepoPath, reqDTO.Branch)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if b {
		return util.AlreadyExistsError()
	}
	_, err = branchmd.InsertProtectedBranch(ctx, reqDTO.RepoPath, reqDTO.Branch)
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
	err := checkAuth(ctx, reqDTO.RepoPath, reqDTO.Operator)
	if err != nil {
		return err
	}
	branch, b, err := branchmd.GetProtectedBranch(ctx, reqDTO.RepoPath, reqDTO.Branch)
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
	err := checkAuth(ctx, reqDTO.RepoPath, reqDTO.Operator)
	if err != nil {
		return ListProtectedBranchRespDTO{}, err
	}
	daoDto := branchmd.ListProtectedBranchReqDTO{
		RepoPath:   reqDTO.RepoPath,
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

// checkAuth 检查权限
func checkAuth(ctx context.Context, repoPath string, operator usermd.UserInfo) error {
	repo, b, err := repomd.GetByPath(ctx, repoPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if !b {
		return util.InvalidArgsError()
	}
	b, err = projectmd.ProjectUserExists(ctx, repo.ProjectId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if !b {
		return util.UnauthorizedError()
	}
	// 不是系统管理员
	if !operator.IsAdmin {
		// 检查是否是仓库管理员
		b, err = repomd.CheckRepoUserExists(ctx, repoPath, operator.Account, repomd.Maintainer)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError()
		}
		if !b {
			return util.UnauthorizedError()
		}
	}
	return nil
}
