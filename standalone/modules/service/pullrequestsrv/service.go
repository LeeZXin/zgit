package pullrequestsrv

import (
	"context"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"path"
	"path/filepath"
	"zgit/pkg/git"
	"zgit/setting"
	"zgit/standalone/modules/model/projectmd"
	"zgit/standalone/modules/model/repomd"
	"zgit/standalone/modules/model/usermd"
	"zgit/util"
)

func PreparePullRequest(ctx context.Context, reqDTO PreparePullRequestReqDTO) (PreparePullRequestRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return PreparePullRequestRespDTO{}, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkPerm(ctx, reqDTO.RepoPath, reqDTO.Operator); err != nil {
		return PreparePullRequestRespDTO{}, err
	}
	absPath := filepath.Join(setting.RepoDir(), reqDTO.RepoPath)
	if !git.CheckRefIsBranch(ctx, absPath, reqDTO.Head) {
		return PreparePullRequestRespDTO{}, util.InvalidArgsError()
	}
	if !git.CheckExists(ctx, absPath, reqDTO.Target) {
		return PreparePullRequestRespDTO{}, util.InvalidArgsError()
	}
	info, err := git.PreparePullRequest(ctx, absPath, reqDTO.Target, reqDTO.Head)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return PreparePullRequestRespDTO{}, util.InternalError()
	}
	ret := PreparePullRequestRespDTO{
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
	return ret, nil
}

func Diff(ctx context.Context, reqDTO DiffReqDTO) (DiffRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return DiffRespDTO{}, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkPerm(ctx, reqDTO.RepoPath, reqDTO.Operator); err != nil {
		return DiffRespDTO{}, err
	}
	absPath := filepath.Join(setting.RepoDir(), reqDTO.RepoPath)
	if !git.CheckExists(ctx, absPath, reqDTO.Target) {
		return DiffRespDTO{}, util.InvalidArgsError()
	}
	if !git.CheckExists(ctx, absPath, reqDTO.Head) {
		return DiffRespDTO{}, util.InvalidArgsError()
	}
	d, err := git.GetDiffFileDetail(ctx, absPath, reqDTO.Target, reqDTO.Head, reqDTO.FileName)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return DiffRespDTO{}, err
	}
	ret := DiffRespDTO{
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

func CatFile(ctx context.Context, reqDTO CatFileReqDTO) ([]DiffLineDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkPerm(ctx, reqDTO.RepoPath, reqDTO.Operator); err != nil {
		return nil, err
	}
	absPath := filepath.Join(setting.RepoDir(), reqDTO.RepoPath)
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

func checkPerm(ctx context.Context, repoPath string, operator usermd.UserInfo) error {
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
		return err
	}
	if !b {
		return util.UnauthorizedError()
	}
	b, err = repomd.CheckRepoUserExists(ctx, repoPath, operator.Account, repomd.Maintainer, repomd.Developer)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if !b {
		return util.InvalidArgsError()
	}
	return nil
}
