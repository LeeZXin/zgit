package repomd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func GenRepoId() string {
	return idutil.RandomUuid()
}

func GetByPath(ctx context.Context, path string) (Repo, bool, error) {
	var ret Repo
	b, err := xormutil.MustGetXormSession(ctx).Where("path = ?", path).Get(&ret)
	return ret, b, err
}

func GetByRepoId(ctx context.Context, repoId string) (Repo, bool, error) {
	var ret Repo
	b, err := xormutil.MustGetXormSession(ctx).Where("repo_id = ?", repoId).Get(&ret)
	return ret, b, err
}

func UpdateTotalAndGitSize(ctx context.Context, repoId string, totalSize, gitSize int64) error {
	_, err := xormutil.MustGetXormSession(ctx).Where("repo_id = ?", repoId).
		Cols("total_size", "git_size").
		Limit(1).
		Update(&Repo{
			TotalSize: totalSize,
			GitSize:   gitSize,
		})
	return err
}

func UpdateIsEmpty(ctx context.Context, repoId string, isEmpty bool) error {
	_, err := xormutil.MustGetXormSession(ctx).Where("repo_id = ?", repoId).
		Cols("is_empty").
		Limit(1).
		Update(&Repo{
			IsEmpty: isEmpty,
		})
	return err
}

func InsertRepo(ctx context.Context, reqDTO InsertRepoReqDTO) (Repo, error) {
	r := Repo{
		RepoId:        GenRepoId(),
		Name:          reqDTO.Name,
		Path:          reqDTO.Path,
		Author:        reqDTO.Author,
		ProjectId:     reqDTO.ProjectId,
		RepoDesc:      reqDTO.RepoDesc,
		DefaultBranch: reqDTO.DefaultBranch,
		RepoType:      int(reqDTO.RepoType),
		IsEmpty:       reqDTO.IsEmpty,
		TotalSize:     reqDTO.TotalSize,
		GitSize:       reqDTO.GitSize,
		LfsSize:       reqDTO.LfsSize,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&r)
	return r, err
}

func DeleteRepo(ctx context.Context, repoId string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("repo_id = ?", repoId).Limit(1).Delete(new(Repo))
	return rows == 1, err
}
