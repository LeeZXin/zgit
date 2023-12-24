package reposrv

import (
	"context"
	"errors"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"path"
	"path/filepath"
	"zgit/pkg/apicode"
	"zgit/pkg/git"
	"zgit/pkg/i18n"
	"zgit/setting"
	"zgit/standalone/modules/model/projectmd"
	"zgit/standalone/modules/model/repomd"
	"zgit/standalone/modules/service/projectsrv"
	"zgit/util"
)

const (
	LsTreeLimit = 25
)

// GetRepoInfoByPath 通过相对路径获取仓库信息
func GetRepoInfoByPath(ctx context.Context, path string) (repomd.RepoInfo, bool, error) {
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByPath(ctx, path)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.RepoInfo{}, false, bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	return repo.ToRepoInfo(), b, nil
}

func TreeRepo(ctx context.Context, reqDTO TreeRepoReqDTO) (TreeRepoRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return TreeRepoRespDTO{}, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return TreeRepoRespDTO{}, util.InternalError()
	}
	if !b {
		return TreeRepoRespDTO{}, util.InvalidArgsError()
	}
	exists, err := projectsrv.ProjectUserExists(ctx, repo.ProjectId, reqDTO.Operator.Account)
	if err != nil {
		return TreeRepoRespDTO{}, err
	}
	if !exists {
		return TreeRepoRespDTO{}, util.UnauthorizedError()
	}
	// 空仓库 需要推代码
	if repo.IsEmpty {
		return TreeRepoRespDTO{IsEmpty: true}, nil
	}
	if reqDTO.Dir == "" {
		reqDTO.Dir = "."
	}
	absPath := filepath.Join(setting.RepoDir(), repo.Path)
	commit, err := git.GetFileLastCommit(ctx, absPath, reqDTO.RefName, reqDTO.Dir)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return TreeRepoRespDTO{}, util.InternalError()
	}
	commits, err := git.LsTreeCommit(ctx, absPath, reqDTO.RefName, reqDTO.Dir, 0, LsTreeLimit)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return TreeRepoRespDTO{}, util.InternalError()
	}
	readme, b, err := git.GetFileContentByRef(ctx, absPath, reqDTO.RefName, filepath.Join(reqDTO.Dir, "readme.md"))
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
	}
	if err == nil && !b {
		readme, b, err = git.GetFileContentByRef(ctx, absPath, reqDTO.RefName, filepath.Join(reqDTO.Dir, "README.md"))
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
		}
	}
	return TreeRepoRespDTO{
		ReadmeText: readme,
		RecentCommit: CommitDTO{
			Author:        commit.Author,
			Committer:     commit.Committer,
			AuthoredDate:  commit.AuthorSigTime,
			CommittedDate: commit.CommitSigTime,
			CommitMsg:     commit.CommitMsg,
			CommitId:      commit.Id,
			ShortId:       util.LongCommitId2ShortId(commit.Id),
		},
		Tree: LsRet2TreeDTO(commits, LsTreeLimit),
	}, nil
}

func LsRet2TreeDTO(commits []git.FileCommit, limit int) TreeDTO {
	files, _ := listutil.Map(commits, func(t git.FileCommit) (FileDTO, error) {
		ret := FileDTO{
			Mode:    t.Mode.Readable(),
			RawPath: t.Path,
			Path:    path.Base(t.Path),
		}
		if t.Commit != nil {
			ret.Commit = CommitDTO{
				Author:        t.Author,
				Committer:     t.Committer,
				AuthoredDate:  t.AuthorSigTime,
				CommittedDate: t.CommitSigTime,
				CommitMsg:     t.CommitMsg,
				CommitId:      t.Blob,
				ShortId:       util.LongCommitId2ShortId(t.Blob),
			}
		}
		return ret, nil
	})
	return TreeDTO{
		Files:   files,
		Limit:   limit,
		HasMore: len(commits) == limit,
	}
}

// InitRepo 初始化仓库
func InitRepo(ctx context.Context, reqDTO InitRepoReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验项目信息
	_, b, err := projectmd.GetByProjectId(ctx, reqDTO.ProjectId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	if !b {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.ProjectNotFound))
	}
	// 相对路径
	relativePath := util.JoinRelativeRepoPath(setting.StandaloneCorpId(), reqDTO.Name)
	// 拼接绝对路径
	absPath := util.JoinAbsRepoPath(setting.StandaloneCorpId(), reqDTO.Name)
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
			Name:          reqDTO.Name,
			Path:          relativePath,
			Author:        reqDTO.Operator.Account,
			ProjectId:     reqDTO.ProjectId,
			RepoDesc:      reqDTO.Desc,
			DefaultBranch: reqDTO.DefaultBranch,
			RepoType:      reqDTO.RepoType,
			IsEmpty:       !reqDTO.CreateReadme && reqDTO.GitIgnoreName == "",
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
				Email:   reqDTO.Operator.Email,
			},
			RepoName:      reqDTO.Name,
			RepoPath:      absPath,
			CreateReadme:  reqDTO.CreateReadme,
			GitIgnoreName: reqDTO.GitIgnoreName,
			DefaultBranch: reqDTO.DefaultBranch,
		})
		if err != nil {
			return err
		}
		// 如果仓库不为空 计算一下仓库大小 lfsSize肯定为0
		if !insertReq.IsEmpty {
			size, err := git.GetRepoSize(absPath)
			logger.Logger.WithContext(ctx).Infof("repo size: %d", size)
			if err == nil {
				repomd.UpdateTotalAndGitSize(ctx, repo.RepoId, size, size)
			}
		}
		return nil
	}); err != nil {
		// 如果有异常 删掉这个仓库
		util.RemoveAll(absPath)
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
		// 是创建人或是管理员
		if reqDTO.Operator.Account == repo.Author || reqDTO.Operator.IsAdmin {
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
