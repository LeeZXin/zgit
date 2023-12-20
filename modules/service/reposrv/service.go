package reposrv

import (
	"context"
	"errors"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"path/filepath"
	"zgit/modules/model/corpmd"
	"zgit/modules/model/projectmd"
	"zgit/modules/model/repomd"
	"zgit/pkg/apicode"
	"zgit/pkg/git"
	"zgit/pkg/i18n"
	"zgit/setting"
	"zgit/util"
)

// GetRepoInfoByPath 通过相对路径获取仓库信息
func GetRepoInfoByPath(ctx context.Context, path string) (repomd.RepoInfo, bool, error) {
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByPath(ctx, path)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.RepoInfo{}, false, errors.New(i18n.GetByKey(i18n.SystemInternalError))
	}
	return repo.ToRepoInfo(), b, nil
}

// InitRepo 初始化仓库
func InitRepo(ctx context.Context, reqDTO InitRepoReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验项目信息
	project, b, err := projectmd.GetByProjectId(ctx, reqDTO.ProjectId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	if !b {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.ProjectNotFound))
	}
	// 企业id对不上
	if project.CorpId != reqDTO.Operator.CorpId {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SystemUnauthorized))
	}
	// 获取企业信息
	corp, b, err := corpmd.GetByCorpId(ctx, reqDTO.Operator.CorpId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	if !b {
		return bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	// 仓库数量大于上限 不给添加
	if corp.RepoLimit >= corp.RepoCount {
		return bizerr.NewBizErr(apicode.OutOfLimitCode.Int(), i18n.GetByKey(i18n.RepoCountOutOfLimit))
	}
	// 相对路径
	relativePath := util.JoinRelativeRepoPath(reqDTO.Operator.CorpId, corp.NodeId, reqDTO.RepoName)
	// 拼接绝对路径
	abPath := util.JoinAbsRepoPath(reqDTO.Operator.CorpId, corp.NodeId, reqDTO.RepoName)
	_, b, err = repomd.GetByPath(ctx, relativePath)
	// 数据库异常
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return err
	}
	// 仓库已存在 不能添加
	if b {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.RepoAlreadyExists))
	}
	// 默认分支
	if reqDTO.DefaultBranch == "" {
		reqDTO.DefaultBranch = setting.DefaultBranch()
	}
	// 数据库事务
	if err = mysqlstore.WithTx(ctx, func(ctx context.Context) error {
		insertReq := repomd.InsertRepoReqDTO{
			Name:          reqDTO.RepoName,
			Path:          relativePath,
			UserId:        reqDTO.Operator.UserId,
			NodeId:        corp.NodeId,
			CorpId:        reqDTO.Operator.CorpId,
			ProjectId:     reqDTO.ProjectId,
			RepoDesc:      reqDTO.RepoDesc,
			DefaultBranch: reqDTO.DefaultBranch,
			RepoType:      reqDTO.RepoType,
			IsEmpty:       reqDTO.CreateReadme || reqDTO.GitIgnoreName != "",
		}
		// 先插入数据库
		repo, err := repomd.InsertRepo(ctx, insertReq)
		if err != nil {
			return err
		}
		// 调用git命令
		err = git.InitRepository(ctx, git.InitRepoOpts{
			Owner: git.User{
				Account: reqDTO.Operator.Account,
				Name:    reqDTO.Operator.Name,
				Email:   reqDTO.Operator.Email,
			},
			RepoName:      reqDTO.RepoName,
			RepoPath:      abPath,
			CreateReadme:  reqDTO.CreateReadme,
			GitIgnoreName: reqDTO.GitIgnoreName,
			DefaultBranch: reqDTO.DefaultBranch,
		})
		if err != nil {
			return err
		}
		// 如果仓库不为空 计算一下仓库大小 lfsSize肯定为0
		if !insertReq.IsEmpty {
			size, err := git.GetRepoSize(abPath)
			if err == nil {
				repomd.UpdateTotalAndGitSize(ctx, repo.RepoId, size, size)
			}
		}
		// 企业仓库计数 +1
		corpmd.IncrRepoCount(ctx, reqDTO.Operator.CorpId)
		return nil
	}); err != nil {
		// 如果有异常 删掉这个仓库
		util.RemoveAll(abPath)
		logger.Logger.WithContext(ctx).Error(err)
		return bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.RepoInitFail))
	}
	return nil
}

// AllGitIgnoreTemplateList 所有gitignore模版名称
func AllGitIgnoreTemplateList() []string {
	return gitignoreSet.AllKeys()
}

// DeleteRepo 删除仓库
func DeleteRepo(ctx context.Context, reqDTO DeleteRepoReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if err := mysqlstore.WithTx(ctx, func(ctx context.Context) error {
		repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
		if err != nil {
			return err
		}
		// 仓库id不存在
		if !b {
			return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.RepoNotFound))
		}
		// 仓库对应的公司id对应不上
		if repo.CorpId != reqDTO.Operator.CorpId {
			return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SystemInvalidArgs))
		}
		// 是创建人或是管理员
		if reqDTO.Operator.UserId == repo.UserId || reqDTO.Operator.IsAdmin {
			// 拼接绝对路径
			absPath := filepath.Join(setting.RepoDir(), repo.Path)
			logger.Logger.WithContext(ctx).Infof("user: %s delete repo: %s", reqDTO.Operator.Account, absPath)
			_, err = repomd.DeleteRepo(ctx, reqDTO.RepoId)
			if err != nil {
				return err
			}
			err = util.RemoveAll(absPath)
			if err != nil {
				return err
			}
			// 企业仓库计数 -1
			corpmd.DeleteCorp(ctx, reqDTO.Operator.CorpId)
			// todo 删除wiki
			return nil
		} else {
			return bizerr.NewBizErr(apicode.UnauthorizedCode.Int(), i18n.GetByKey(i18n.SystemUnauthorized))
		}
	}); err != nil {
		if _, ok := err.(*bizerr.Err); ok {
			return err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return errors.New(i18n.GetByKey(i18n.SystemInternalError))
	}
	return nil
}
