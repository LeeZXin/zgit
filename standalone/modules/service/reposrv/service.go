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
	"zgit/standalone/modules/model/usermd"
	"zgit/util"
)

const (
	LsTreeLimit = 25
)

// GetInfoByPath 通过相对路径获取仓库信息
func GetInfoByPath(ctx context.Context, path string) (repomd.RepoInfo, bool, error) {
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByPath(ctx, path)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.RepoInfo{}, false, bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	return repo.ToRepoInfo(), b, nil
}

func EntriesRepo(ctx context.Context, reqDTO EntriesRepoReqDTO) (TreeDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return TreeDTO{}, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	repo, err := checkAuth(ctx, reqDTO.RepoPath, reqDTO.Operator)
	if err != nil {
		return TreeDTO{}, err
	}
	// 空仓库 需要推代码
	if repo.IsEmpty {
		return TreeDTO{}, nil
	}
	if reqDTO.Dir == "" {
		reqDTO.Dir = "."
	}
	absPath := filepath.Join(setting.RepoDir(), repo.Path)
	commits, err := git.LsTreeCommit(ctx, absPath, reqDTO.RefName, reqDTO.Dir, reqDTO.Offset, LsTreeLimit)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return TreeDTO{}, util.InternalError()
	}
	return lsRet2TreeDTO(commits, reqDTO.Offset, LsTreeLimit), nil
}

// ListRepo 展示仓库列表
func ListRepo(ctx context.Context, reqDTO ListRepoReqDTO) (ListRepoRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return ListRepoRespDTO{}, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	b, err := projectmd.ProjectUserExists(ctx, reqDTO.ProjectId, reqDTO.Operator.Account)
	if err != nil {
		return ListRepoRespDTO{}, err
	}
	if !b {
		return ListRepoRespDTO{}, util.UnauthorizedError()
	}
	ret := ListRepoRespDTO{
		Limit: reqDTO.Limit,
	}
	ret.TotalCount, err = repomd.CountRepo(ctx, reqDTO.ProjectId, reqDTO.SearchName)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return ListRepoRespDTO{}, util.InternalError()
	}
	if ret.TotalCount > 0 {
		ret.RepoList, err = repomd.ListRepo(ctx, reqDTO.Offset, reqDTO.Limit, reqDTO.ProjectId, reqDTO.SearchName)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return ListRepoRespDTO{}, util.InternalError()
		}
		if len(ret.RepoList) > 0 {
			ret.Cursor = ret.RepoList[len(ret.RepoList)-1].Id
		}
	}
	return ret, nil
}

// CatFile 展示文件内容
func CatFile(ctx context.Context, reqDTO CatFileReqDTO) (git.FileMode, string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", "", err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	if _, err := checkAuth(ctx, reqDTO.RepoPath, reqDTO.Operator); err != nil {
		return "", "", err
	}
	absPath := filepath.Join(setting.RepoDir(), reqDTO.RepoPath)
	fileMode, content, _, err := git.GetFileContentByRef(ctx, absPath, reqDTO.RefName, reqDTO.FileName)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", "", util.InternalError()
	}
	return fileMode, content, nil
}

// TreeRepo 代码基本数据
func TreeRepo(ctx context.Context, reqDTO TreeRepoReqDTO) (TreeRepoRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return TreeRepoRespDTO{}, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	repo, err := checkAuth(ctx, reqDTO.RepoPath, reqDTO.Operator)
	if err != nil {
		return TreeRepoRespDTO{}, err
	}
	// 空仓库 需要推代码
	if repo.IsEmpty {
		return TreeRepoRespDTO{IsEmpty: true}, nil
	}
	if reqDTO.Dir == "" {
		reqDTO.Dir = "."
	}
	absPath := filepath.Join(setting.RepoDir(), reqDTO.RepoPath)
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
	_, readme, hasReadme, err := git.GetFileContentByRef(ctx, absPath, reqDTO.RefName, filepath.Join(reqDTO.Dir, "readme.md"))
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
	}
	if err == nil && !hasReadme {
		_, readme, hasReadme, err = git.GetFileContentByRef(ctx, absPath, reqDTO.RefName, filepath.Join(reqDTO.Dir, "README.md"))
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
		}
	}
	return TreeRepoRespDTO{
		ReadmeText: readme,
		HasReadme:  hasReadme,
		RecentCommit: CommitDTO{
			Author:        commit.Author,
			Committer:     commit.Committer,
			AuthoredDate:  commit.AuthorSigTime,
			CommittedDate: commit.CommitSigTime,
			CommitMsg:     commit.CommitMsg,
			CommitId:      commit.Id,
			ShortId:       util.LongCommitId2ShortId(commit.Id),
		},
		Tree: lsRet2TreeDTO(commits, 0, LsTreeLimit),
	}, nil
}

func lsRet2TreeDTO(commits []git.FileCommit, offset, limit int) TreeDTO {
	files, _ := listutil.Map(commits, func(t git.FileCommit) (FileDTO, error) {
		ret := FileDTO{
			Mode:    t.Mode.Readable(),
			RawPath: t.Path,
			Path:    path.Base(t.Path),
			Commit: CommitDTO{
				Author:        t.Author,
				Committer:     t.Committer,
				AuthoredDate:  t.AuthorSigTime,
				CommittedDate: t.CommitSigTime,
				CommitMsg:     t.CommitMsg,
				CommitId:      t.Blob,
				ShortId:       util.LongCommitId2ShortId(t.Blob),
			},
		}
		return ret, nil
	})
	return TreeDTO{
		Files:   files,
		Limit:   limit,
		Offset:  offset,
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
	b, err := projectmd.ProjectUserExists(ctx, reqDTO.ProjectId, reqDTO.Operator.Account)
	if err != nil {
		return err
	}
	if !b {
		return util.UnauthorizedError()
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
				repomd.UpdateTotalAndGitSize(ctx, repo.Path, size, size)
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

// AllTypeList 所有仓库类型
func AllTypeList() []RepoTypeDTO {
	return []RepoTypeDTO{
		{
			Option: repomd.InternalRepoType.Int(),
			Name:   repomd.InternalRepoType.Readable(),
		},
		{
			Option: repomd.PublicRepoType.Int(),
			Name:   repomd.PublicRepoType.Readable(),
		},
		{
			Option: repomd.PrivateRepoType.Int(),
			Name:   repomd.PrivateRepoType.Readable(),
		},
	}
}

// DeleteRepo 删除仓库
func DeleteRepo(ctx context.Context, reqDTO DeleteRepoReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	repo, err := checkAuth(ctx, reqDTO.RepoPath, reqDTO.Operator)
	if err != nil {
		return err
	}
	// 如果不是系统管理员 检查仓库管理员权限
	if !reqDTO.Operator.IsAdmin {
		b, err := repomd.CheckRepoUserExists(ctx, reqDTO.RepoPath, reqDTO.Operator.Account, repomd.Maintainer)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError()
		}
		if !b {
			return util.UnauthorizedError()
		}
	}
	// 拼接绝对路径
	absPath := filepath.Join(setting.RepoDir(), reqDTO.RepoPath)
	logger.Logger.WithContext(ctx).Infof("user: %s delete repo: %s", reqDTO.Operator.Account, absPath)
	if err := mysqlstore.WithTx(ctx, func(ctx context.Context) error {
		_, err := repomd.DeleteRepo(ctx, repo)
		if err != nil {
			return err
		}
		err = util.RemoveAll(absPath)
		if err != nil {
			return err
		}
		// todo 删除wiki
		return nil
	}); err != nil {
		if _, ok := err.(*bizerr.Err); ok {
			return err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return errors.New(i18n.GetByKey(i18n.SystemInternalError))
	}
	return nil
}

// AllBranches 仓库所有分支
func AllBranches(ctx context.Context, reqDTO AllBranchesReqDTO) ([]string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if _, err := checkAuth(ctx, reqDTO.RepoPath, reqDTO.Operator); err != nil {
		return nil, err
	}
	absPath := filepath.Join(setting.RepoDir(), reqDTO.RepoPath)
	branchList, err := git.GetAllBranchList(ctx, absPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError()
	}
	return branchList, nil
}

// AllTags 仓库所有tag
func AllTags(ctx context.Context, reqDTO AllTagsReqDTO) ([]string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if _, err := checkAuth(ctx, reqDTO.RepoPath, reqDTO.Operator); err != nil {
		return nil, err
	}
	absPath := filepath.Join(setting.RepoDir(), reqDTO.RepoPath)
	tagList, err := git.GetAllTagList(ctx, absPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError()
	}
	return tagList, nil
}

func checkAuth(ctx context.Context, repoPath string, operator usermd.UserInfo) (repomd.Repo, error) {
	repo, b, err := repomd.GetByPath(ctx, repoPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, util.InternalError()
	}
	if !b {
		return repomd.Repo{}, util.InvalidArgsError()
	}
	b, err = projectmd.ProjectUserExists(ctx, repo.ProjectId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, util.InternalError()
	}
	if !b {
		return repomd.Repo{}, util.UnauthorizedError()
	}
	b, err = repomd.CheckRepoUserExists(ctx, repoPath, operator.Account, repomd.Guest, repomd.Maintainer, repomd.Developer)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, util.InternalError()
	}
	if !b {
		return repomd.Repo{}, util.UnauthorizedError()
	}
	return repo, nil
}
