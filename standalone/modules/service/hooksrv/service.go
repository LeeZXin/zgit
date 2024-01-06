package hooksrv

import (
	"context"
	"github.com/IGLOU-EU/go-wildcard/v2"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"path/filepath"
	"strings"
	"zgit/pkg/apicode"
	"zgit/pkg/git"
	"zgit/pkg/hook"
	"zgit/pkg/i18n"
	"zgit/setting"
	"zgit/standalone/modules/model/branchmd"
	"zgit/standalone/modules/model/repomd"
	"zgit/util"
)

func PreReceive(ctx context.Context, opts hook.Opts) error {
	logger.Logger.WithContext(ctx).Info("pre-receive", opts)
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, opts.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if !b {
		return util.InvalidArgsError()
	}
	repoPath := filepath.Join(setting.RepoDir(), repo.Path)
	var pbList []branchmd.ProtectedBranchDTO
	for _, info := range opts.RevInfoList {
		name := info.RefName
		// 是分支
		if strings.HasPrefix(name, git.BranchPrefix) {
			// 检查是否是保护分支
			if pbList == nil {
				// 懒加载一下
				pbList, err = branchmd.ListProtectedBranch(ctx, opts.RepoId)
				if err != nil {
					logger.Logger.WithContext(ctx).Error(err)
					return util.InternalError()
				}
			}
			name = strings.TrimPrefix(name, git.BranchPrefix)
			for _, pb := range pbList {
				// 通配符匹配 是保护分支
				if wildcard.Match(pb.Branch, name) {
					// 只有可推送名单里面才能直接push
					if opts.PrId == "" && len(pb.Cfg.DirectPushList) > 0 {
						// prId为空说明不是来自合并请求的push
						contains, _ := listutil.Contains(pb.Cfg.DirectPushList, func(account string) (bool, error) {
							return account == opts.PusherId, nil
						})
						if !contains {
							return util.NewBizErr(apicode.ForcePushForbiddenCode, i18n.ProtectedBranchNotAllowDirectPush)
						}
					}
					// 不允许删除保护分支
					if info.NewCommitId == git.ZeroCommitId {
						return util.NewBizErr(apicode.ForcePushForbiddenCode, i18n.ProtectedBranchNotAllowDelete)
					}
					// 检查push -f
					isForcePush, err := git.DetectForcePush(ctx,
						repoPath,
						info.OldCommitId,
						info.NewCommitId,
						git.DetectForcePushEnv{
							ObjectDirectory:              opts.ObjectDirectory,
							AlternativeObjectDirectories: opts.AlternativeObjectDirectories,
							QuarantinePath:               opts.QuarantinePath,
						})
					if err != nil {
						logger.Logger.WithContext(ctx).Error(err)
						return util.InternalError()
					}
					if isForcePush {
						// 禁止push -f
						return util.NewBizErr(apicode.ForcePushForbiddenCode, i18n.ProtectedBranchNotAllowForcePush)
					}
				}
			}
		}
	}
	return nil
}

func PostReceive(ctx context.Context, opts hook.Opts) error {
	logger.Logger.WithContext(ctx).Info("post-receive", opts)
	return nil
}
