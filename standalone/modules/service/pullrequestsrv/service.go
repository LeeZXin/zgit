package pullrequestsrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"path/filepath"
	"zgit/pkg/apicode"
	"zgit/pkg/git"
	"zgit/pkg/i18n"
	"zgit/pkg/perm"
	"zgit/setting"
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
	repo, p, err := getPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return err
	}
	if !p.CanHandlePullRequest {
		return util.UnauthorizedError()
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
	// 可合并提交为0
	if len(info.Commits) == 0 || len(info.ConflictFiles) > 0 {
		return bizerr.NewBizErr(apicode.PullRequestCannotMergeCode.Int(), i18n.GetByKey(i18n.PullRequestCannotMerge))
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
	pr, b, err := pullrequestmd.GetByPrId(ctx, reqDTO.PrId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 只允许从open -> closed
	if pr.PrStatus != pullrequestmd.PrOpenStatus.Int() {
		return util.InvalidArgsError()
	}
	// 校验权限
	_, p, err := getPerm(ctx, pr.RepoId, reqDTO.Operator)
	if err != nil {
		return err
	}
	if !p.CanHandlePullRequest {
		return util.UnauthorizedError()
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
	pr, b, err := pullrequestmd.GetByPrId(ctx, reqDTO.PrId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 只允许从open -> closed
	if pr.PrStatus != pullrequestmd.PrOpenStatus.Int() {
		return util.InvalidArgsError()
	}
	// 校验权限
	repo, p, err := getPerm(ctx, pr.RepoId, reqDTO.Operator)
	if err != nil {
		return err
	}
	if !p.CanHandlePullRequest {
		return util.UnauthorizedError()
	}
	absPath := filepath.Join(setting.RepoDir(), repo.Path)
	info, err := git.GetDiffCommitsInfo(ctx, absPath, pr.Target, pr.Head)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	// 可合并提交为0
	if len(info.Commits) == 0 || len(info.ConflictFiles) > 0 {
		return bizerr.NewBizErr(apicode.PullRequestCannotMergeCode.Int(), i18n.GetByKey(i18n.PullRequestCannotMerge))
	}
	return mysqlstore.WithTx(ctx, func(ctx context.Context) error {
		b, err = pullrequestmd.UpdatePrStatusAndCommitId(
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

func getPerm(ctx context.Context, repoId string, operator usermd.UserInfo) (repomd.Repo, perm.RepoPerm, error) {
	repo, b, err := repomd.GetByRepoId(ctx, repoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, perm.RepoPerm{}, util.InternalError()
	}
	if !b {
		return repomd.Repo{}, perm.RepoPerm{}, util.InvalidArgsError()
	}
	p, b, err := projectmd.GetProjectUserPermDetail(ctx, repo.ProjectId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, perm.RepoPerm{}, util.InternalError()
	}
	if !b {
		return repomd.Repo{}, perm.RepoPerm{}, util.UnauthorizedError()
	}
	return repo, p.PermDetail.GetRepoPerm(repoId), nil
}
