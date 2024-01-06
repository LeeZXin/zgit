package branchsrv

import (
	"context"
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
	repo, err := checkPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
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
		// 评审账号不合法
		if !b {
			return util.NewBizErr(apicode.InvalidArgsCode, i18n.UserAccountNotFoundWarnFormat, account)
		}
		// 检查评审者是否有访问代码的权限
		detail, b, err := projectmd.GetProjectUserPermDetail(ctx, repo.ProjectId, account)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError()
		}
		if !b || !detail.PermDetail.GetRepoPerm(repo.RepoId).CanAccess {
			return util.NewBizErr(apicode.InvalidArgsCode, i18n.UserAccountUnauthorizedReviewCodeWarnFormat, account)
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
	pb, b, err := branchmd.GetProtectedBranchByBid(ctx, reqDTO.Bid)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if !b {
		return util.InvalidArgsError()
	}
	_, err = checkPerm(ctx, pb.RepoId, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = branchmd.DeleteProtectedBranch(ctx, pb)
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
	_, err := checkPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
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
			Bid:    t.Bid,
			RepoId: t.RepoId,
			Branch: t.Branch,
			Cfg:    t.Cfg,
		}, nil
	})
	return ret, nil
}

// checkPerm 检查权限 只需检查是否是项目管理员
func checkPerm(ctx context.Context, repoId string, operator usermd.UserInfo) (repomd.Repo, error) {
	// 检查仓库是否存在
	repo, b, err := repomd.GetByRepoId(ctx, repoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, util.InternalError()
	}
	if !b {
		return repomd.Repo{}, util.InvalidArgsError()
	}
	// 如果不是 检查用户组权限
	permDetail, b, err := projectmd.GetProjectUserPermDetail(ctx, repo.ProjectId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repo, util.InternalError()
	}
	if !b {
		return repo, util.UnauthorizedError()
	}
	if !permDetail.PermDetail.GetRepoPerm(repoId).CanHandleProtectedBranch {
		return repo, util.UnauthorizedError()
	}
	return repo, nil
}
