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
	"zgit/pkg/perm"
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
		return repomd.RepoInfo{}, false, util.InternalError()
	}
	return repo.ToRepoInfo(), b, nil
}

func EntriesRepo(ctx context.Context, reqDTO EntriesRepoReqDTO) (TreeDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return TreeDTO{}, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	repo, p, err := getPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return TreeDTO{}, err
	}
	if p.GetRepoPerm(repo.RepoId).CanAccess {
		return TreeDTO{}, util.UnauthorizedError()
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
func ListRepo(ctx context.Context, reqDTO ListRepoReqDTO) ([]repomd.Repo, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	p, b, err := projectmd.GetProjectUserPermDetail(ctx, reqDTO.ProjectId, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError()
	}
	if !b {
		return nil, util.UnauthorizedError()
	}
	// 项目管理员可看到所有仓库或者应用所有仓库权限配置
	if p.IsAdmin || p.PermDetail.ApplyDefaultRepoPerm {
		repoList, err := repomd.ListAllRepo(ctx, reqDTO.ProjectId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError()
		}
		return repoList, nil
	}
	// 通过可访问仓库id查询
	permList := p.PermDetail.RepoPermList
	repoIdList, _ := listutil.Map(permList, func(t perm.RepoPermWithId) (string, error) {
		return t.RepoId, nil
	})
	repoList, err := repomd.ListRepoByIdList(ctx, reqDTO.ProjectId, repoIdList)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError()
	}
	return repoList, nil
}

// CatFile 展示文件内容
func CatFile(ctx context.Context, reqDTO CatFileReqDTO) (git.FileMode, string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", "", err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	repo, p, err := getPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return "", "", err
	}
	if !p.GetRepoPerm(repo.RepoId).CanAccess {
		return "", "", util.UnauthorizedError()
	}
	absPath := filepath.Join(setting.RepoDir(), repo.Path)
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
	repo, p, err := getPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return TreeRepoRespDTO{}, err
	}
	if !p.GetRepoPerm(repo.RepoId).CanAccess {
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
	p, b, err := projectmd.GetProjectUserPermDetail(ctx, reqDTO.ProjectId, reqDTO.Operator.Account)
	if err != nil {
		return err
	}
	if !b {
		return util.UnauthorizedError()
	}
	// 是否可创建项目
	if p.PermDetail.ProjectPerm.CanInitRepo {
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
		return util.NewBizErr(apicode.InvalidArgsCode, i18n.RepoAlreadyExists)
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
		return util.InternalError()
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
	repo, p, err := getPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return err
	}
	// 是否可删除权限
	if !p.ProjectPerm.CanDeleteRepo {
		return util.UnauthorizedError()
	}
	// 拼接绝对路径
	absPath := filepath.Join(setting.RepoDir(), repo.Path)
	logger.Logger.WithContext(ctx).Infof("user: %s delete repo: %s", reqDTO.Operator.Account, absPath)
	if err = mysqlstore.WithTx(ctx, func(ctx context.Context) error {
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
	repo, p, err := getPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	// 是否可访问
	if !p.GetRepoPerm(repo.RepoId).CanAccess {
		return nil, util.UnauthorizedError()
	}
	absPath := filepath.Join(setting.RepoDir(), repo.Path)
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
	repo, p, err := getPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	// 是否可访问
	if !p.GetRepoPerm(repo.RepoId).CanAccess {
		return nil, util.UnauthorizedError()
	}
	absPath := filepath.Join(setting.RepoDir(), repo.Path)
	tagList, err := git.GetAllTagList(ctx, absPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError()
	}
	return tagList, nil
}

// Gc git gc
func Gc(ctx context.Context, reqDTO GcReqDTO) error {
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
	// 是否可访问
	if !p.GetRepoPerm(repo.RepoId).CanAccess {
		return util.UnauthorizedError()
	}
	absPath := filepath.Join(setting.RepoDir(), repo.Path)
	if err = git.Gc(ctx, absPath); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	size, err := git.GetRepoSize(absPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	totalSize := repo.LfsSize + repo.WikiSize + size
	if err = repomd.UpdateTotalAndGitSize(ctx, reqDTO.RepoId, totalSize, size); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	return nil
}

func DiffCommits(ctx context.Context, reqDTO DiffCommitsReqDTO) (DiffCommitsRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return DiffCommitsRespDTO{}, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	repo, p, err := getPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return DiffCommitsRespDTO{}, err
	}
	if !p.GetRepoPerm(repo.RepoId).CanAccess {
		return DiffCommitsRespDTO{}, util.UnauthorizedError()
	}
	absPath := filepath.Join(setting.RepoDir(), repo.Path)
	if !git.CheckExists(ctx, absPath, reqDTO.Head) {
		return DiffCommitsRespDTO{}, util.InvalidArgsError()
	}
	if !git.CheckExists(ctx, absPath, reqDTO.Target) {
		return DiffCommitsRespDTO{}, util.InvalidArgsError()
	}
	info, err := git.GetDiffCommitsInfo(ctx, absPath, reqDTO.Target, reqDTO.Head)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return DiffCommitsRespDTO{}, util.InternalError()
	}
	ret := DiffCommitsRespDTO{
		Target:       info.Target,
		Head:         info.Head,
		TargetCommit: commit2Dto(info.TargetCommit),
		HeadCommit:   commit2Dto(info.HeadCommit),
		NumFiles:     info.NumFiles,
		DiffNumsStats: DiffNumsStatInfoDTO{
			FileChangeNums: info.DiffNumsStats.FileChangeNums,
			InsertNums:     info.DiffNumsStats.InsertNums,
			DeleteNums:     info.DiffNumsStats.DeleteNums,
		},
		ConflictFiles: info.ConflictFiles,
	}
	ret.DiffNumsStats.Stats, _ = listutil.Map(info.DiffNumsStats.Stats, func(t git.DiffNumsStat) (DiffNumsStatDTO, error) {
		return DiffNumsStatDTO{
			RawPath:    t.Path,
			Path:       path.Base(t.Path),
			TotalNums:  t.TotalNums,
			InsertNums: t.InsertNums,
			DeleteNums: t.DeleteNums,
		}, nil
	})
	ret.Commits, _ = listutil.Map(info.Commits, func(t git.Commit) (CommitDTO, error) {
		return commit2Dto(t), nil
	})
	ret.CanMerge = len(ret.Commits) > 0 && len(ret.ConflictFiles) == 0
	return ret, nil
}

func DiffFile(ctx context.Context, reqDTO DiffFileReqDTO) (DiffFileRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return DiffFileRespDTO{}, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	repo, p, err := getPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return DiffFileRespDTO{}, err
	}
	if !p.GetRepoPerm(repo.RepoId).CanAccess {
		return DiffFileRespDTO{}, util.UnauthorizedError()
	}
	absPath := filepath.Join(setting.RepoDir(), repo.Path)
	if !git.CheckExists(ctx, absPath, reqDTO.Target) {
		return DiffFileRespDTO{}, util.InvalidArgsError()
	}
	if !git.CheckExists(ctx, absPath, reqDTO.Head) {
		return DiffFileRespDTO{}, util.InvalidArgsError()
	}
	d, err := git.GetDiffFileDetail(ctx, absPath, reqDTO.Target, reqDTO.Head, reqDTO.FileName)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return DiffFileRespDTO{}, err
	}
	ret := DiffFileRespDTO{
		FilePath:    d.FilePath,
		OldMode:     d.OldMode,
		Mode:        d.Mode,
		IsSubModule: d.IsSubModule,
		FileType:    d.FileType,
		IsBinary:    d.IsBinary,
		RenameFrom:  d.RenameFrom,
		RenameTo:    d.RenameTo,
		CopyFrom:    d.CopyFrom,
		CopyTo:      d.CopyTo,
	}
	ret.Lines, _ = listutil.Map(d.Lines, func(t git.DiffLine) (DiffLineDTO, error) {
		return DiffLineDTO{
			Index:   t.Index,
			LeftNo:  t.LeftNo,
			Prefix:  t.Prefix,
			RightNo: t.RightNo,
			Text:    t.Text,
		}, nil
	})
	return ret, nil
}

func commit2Dto(commit git.Commit) CommitDTO {
	return CommitDTO{
		Author:        commit.Author,
		Committer:     commit.Committer,
		AuthoredDate:  commit.AuthorSigTime,
		CommittedDate: commit.CommitSigTime,
		CommitMsg:     commit.CommitMsg,
		CommitId:      commit.Id,
		ShortId:       util.LongCommitId2ShortId(commit.Id),
	}
}

func ShowDiffTextContent(ctx context.Context, reqDTO ShowDiffTextContentReqDTO) ([]DiffLineDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	repo, p, err := getPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	if !p.GetRepoPerm(repo.RepoId).CanAccess {
		return nil, util.UnauthorizedError()
	}
	absPath := filepath.Join(setting.RepoDir(), repo.Path)
	if !git.CheckRefIsCommit(ctx, absPath, reqDTO.CommitId) {
		return nil, util.InvalidArgsError()
	}
	var startLine int
	if reqDTO.Direction == UpDirection {
		if reqDTO.Limit < 0 {
			startLine = 0
		} else {
			startLine = reqDTO.Offset - reqDTO.Limit
		}
	} else {
		startLine = reqDTO.Offset
	}
	if startLine < 0 {
		startLine = 0
	}
	lineList, err := git.ShowFileTextContentByCommitId(ctx, absPath, reqDTO.CommitId, reqDTO.FileName, startLine, reqDTO.Limit)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError()
	}
	ret := make([]DiffLineDTO, 0, len(lineList))
	for i, line := range lineList {
		n := startLine + i
		ret = append(ret, DiffLineDTO{
			Index:   i,
			LeftNo:  n,
			Prefix:  git.NormalLinePrefix,
			RightNo: n,
			Text:    line,
		})
	}
	return ret, nil
}

func getPerm(ctx context.Context, repoId string, operator usermd.UserInfo) (repomd.Repo, perm.Detail, error) {
	repo, b, err := repomd.GetByRepoId(ctx, repoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, perm.Detail{}, util.InternalError()
	}
	if !b {
		return repomd.Repo{}, perm.Detail{}, util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return repo, perm.DefaultPermDetail, nil
	}
	p, b, err := projectmd.GetProjectUserPermDetail(ctx, repo.ProjectId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, perm.Detail{}, util.InternalError()
	}
	if !b {
		return repomd.Repo{}, perm.Detail{}, util.UnauthorizedError()
	}
	return repo, p.PermDetail, nil
}
