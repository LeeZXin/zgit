package pullrequestsrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"path/filepath"
	"zgit/pkg/apicode"
	"zgit/pkg/git"
	"zgit/pkg/i18n"
	"zgit/setting"
	"zgit/standalone/modules/model/branchmd"
	"zgit/standalone/modules/model/projectmd"
	"zgit/standalone/modules/model/pullrequestmd"
	"zgit/standalone/modules/model/repomd"
	"zgit/standalone/modules/model/usermd"
	"zgit/util"
)

func SubmitPullRequest(ctx context.Context, reqDTO SubmitPullRequestReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	repo, err := checkPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return err
	}
	absPath := filepath.Join(setting.RepoDir(), repo.Path)
	if !git.CheckRefIsBranch(ctx, absPath, reqDTO.Head) {
		return util.InvalidArgsError()
	}
	if !git.CheckExists(ctx, absPath, reqDTO.Target) {
		return util.InvalidArgsError()
	}
	info, err := git.GetDiffCommitsInfo(ctx, absPath, reqDTO.Target, reqDTO.Head)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	// 不可合并
	if !info.IsMergeAble() {
		return util.NewBizErr(apicode.PullRequestCannotMergeCode, i18n.PullRequestCannotMerge)
	}
	_, err = pullrequestmd.InsertPullRequest(ctx, pullrequestmd.InsertPullRequestReqDTO{
		RepoId:   reqDTO.RepoId,
		Target:   reqDTO.Target,
		Head:     reqDTO.Head,
		CreateBy: reqDTO.Operator.Account,
		PrStatus: pullrequestmd.PrOpenStatus,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	return nil
}

func ClosePullRequest(ctx context.Context, reqDTO ClosePullRequestReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	pr, _, err := checkPerm(ctx, reqDTO.PrId, reqDTO.Operator)
	if err != nil {
		return err
	}
	// 只允许从open -> closed
	if pr.PrStatus != pullrequestmd.PrOpenStatus {
		return util.InvalidArgsError()
	}
	_, err = pullrequestmd.UpdatePrStatus(ctx, reqDTO.PrId, pullrequestmd.PrOpenStatus, pullrequestmd.PrClosedStatus)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	return nil
}

func MergePullRequest(ctx context.Context, reqDTO MergePullRequestReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	pr, repo, err := checkPerm(ctx, reqDTO.PrId, reqDTO.Operator)
	if err != nil {
		return err
	}
	// 只允许从open -> closed
	if pr.PrStatus != pullrequestmd.PrOpenStatus {
		return util.InvalidArgsError()
	}
	// 检查是否是保护分支
	cfg, isProtectedBranch, err := branchmd.IsProtectedBranch(ctx, pr.RepoId, pr.Head)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if isProtectedBranch {
		// 检查评审配置 评审者数量大于0
		if cfg.ReviewCountWhenCreatePr > 0 {
			reviewCount, err := pullrequestmd.CountReview(ctx, reqDTO.PrId, pullrequestmd.AgreeMergeStatus)
			if err != nil {
				logger.Logger.WithContext(ctx).Error(err)
				return util.InternalError()
			}
			// 小于配置数量 不可合并
			if reviewCount < cfg.ReviewCountWhenCreatePr {
				return util.NewBizErr(apicode.PullRequestCannotMergeCode, i18n.PullRequestReviewerCountLowerThanCfg)
			}
		}
	}
	absPath := filepath.Join(setting.RepoDir(), repo.Path)
	info, err := git.GetDiffCommitsInfo(ctx, absPath, pr.Target, pr.Head)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	// 不可合并
	if !info.IsMergeAble() {
		return util.NewBizErr(apicode.PullRequestCannotMergeCode, i18n.PullRequestCannotMerge)
	}
	return mysqlstore.WithTx(ctx, func(ctx context.Context) error {
		b, err := pullrequestmd.UpdatePrStatusAndCommitId(
			ctx,
			pr.PrId,
			pullrequestmd.PrOpenStatus,
			pullrequestmd.PrMergedStatus,
			info.TargetCommit.Id,
			info.HeadCommit.Id,
		)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError()
		}
		if b {
			err = git.Merge(ctx, absPath, pr.Target, pr.Head, info, git.MergeRepoOpts{
				PrId:    pr.PrId,
				Message: fmt.Sprintf(i18n.GetByKey(i18n.PullRequestMergeMessage), pr.PrId, pr.CreateBy, reqDTO.Operator.Account),
			})
			if err != nil {
				logger.Logger.WithContext(ctx).Error(err)
				return util.InternalError()
			}
		}
		return nil
	})
}

func ReviewPullRequest(ctx context.Context, reqDTO ReviewPullRequestReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	// 检查是否重复提交
	_, b, err := pullrequestmd.GetReview(ctx, reqDTO.PrId, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if b {
		return util.NewBizErr(apicode.DataAlreadyExistsCode, i18n.RepoAlreadyExists)
	}
	// 检查评审者是否有访问代码的权限
	pr, b, err := pullrequestmd.GetByPrId(ctx, reqDTO.PrId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if !b {
		return util.InvalidArgsError()
	}
	repo, b, err := repomd.GetByRepoId(ctx, pr.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 系统管理员有所有权限
	if !reqDTO.Operator.IsAdmin {
		p, b, err := projectmd.GetProjectUserPermDetail(ctx, repo.ProjectId, reqDTO.Operator.Account)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError()
		}
		if !b {
			return util.InvalidArgsError()
		}
		if !p.PermDetail.GetRepoPerm(pr.RepoId).CanAccess {
			return util.UnauthorizedError()
		}
	}
	// 检查是否是保护分支
	cfg, isProtectedBranch, err := branchmd.IsProtectedBranch(ctx, repo.RepoId, pr.Head)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if isProtectedBranch {
		// 看看是否在评审名单里面
		if len(cfg.ReviewerList) > 0 {
			contains, _ := listutil.Contains(cfg.ReviewerList, func(account string) (bool, error) {
				return account == reqDTO.Operator.Account, nil
			})
			if !contains {
				return util.UnauthorizedError()
			}
		}
	}
	err = pullrequestmd.InsertReview(ctx, pullrequestmd.InsertReviewReqDTO{
		PrId:      reqDTO.PrId,
		ReviewMsg: reqDTO.ReviewMsg,
		Status:    reqDTO.Status,
		Reviewer:  reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	return nil
}

// checkPerm 校验权限
func checkPerm(ctx context.Context, prId string, operator usermd.UserInfo) (pullrequestmd.PullRequest, repomd.Repo, error) {
	pr, b, err := pullrequestmd.GetByPrId(ctx, prId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return pullrequestmd.PullRequest{}, repomd.Repo{}, util.InternalError()
	}
	if !b {
		return pullrequestmd.PullRequest{}, repomd.Repo{}, util.InvalidArgsError()
	}
	repo, err := checkPermByRepoId(ctx, pr.RepoId, operator)
	return pr, repo, err
}

// checkPermByRepoId 校验权限
func checkPermByRepoId(ctx context.Context, repoId string, operator usermd.UserInfo) (repomd.Repo, error) {
	repo, b, err := repomd.GetByRepoId(ctx, repoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, util.InternalError()
	}
	if !b {
		return repomd.Repo{}, util.InvalidArgsError()
	}
	// 如果是系统管理员有所有权限
	if operator.IsAdmin {
		return repo, nil
	}
	p, b, err := projectmd.GetProjectUserPermDetail(ctx, repo.ProjectId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repo, util.InternalError()
	}
	if !b {
		return repo, util.InvalidArgsError()
	}
	if !p.PermDetail.GetRepoPerm(repoId).CanHandlePullRequest {
		return repo, util.UnauthorizedError()
	}
	return repo, nil
}
