package pullrequestsrv

import (
	"context"
	"fmt"
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
	"zgit/standalone/modules/model/pullrequestmd"
	"zgit/standalone/modules/model/repomd"
	"zgit/standalone/modules/model/usermd"
	"zgit/util"
)

func DiffCommits(ctx context.Context, reqDTO DiffCommitsReqDTO) (DiffCommitsRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return DiffCommitsRespDTO{}, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	repo, err := checkPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return DiffCommitsRespDTO{}, err
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

func SubmitPullRequest(ctx context.Context, reqDTO SubmitPullRequestReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	repo, err := checkPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
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
	if _, err = checkPerm(ctx, pr.RepoId, reqDTO.Operator); err != nil {
		return err
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
	repo, err := checkPerm(ctx, pr.RepoId, reqDTO.Operator)
	if err != nil {
		return err
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

func DiffFile(ctx context.Context, reqDTO DiffFileReqDTO) (DiffFileRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return DiffFileRespDTO{}, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	repo, err := checkPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return DiffFileRespDTO{}, err
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

func CatFile(ctx context.Context, reqDTO CatFileReqDTO) ([]DiffLineDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	repo, err := checkPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return nil, err
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

func checkPerm(ctx context.Context, repoId string, operator usermd.UserInfo) (repomd.Repo, error) {
	repo, b, err := repomd.GetByRepoId(ctx, repoId)
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
	b, err = repomd.CheckRepoUserExists(ctx, repoId, operator.Account, repomd.Maintainer, repomd.Developer)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, util.InternalError()
	}
	if !b {
		return repomd.Repo{}, util.InvalidArgsError()
	}
	return repo, nil
}
