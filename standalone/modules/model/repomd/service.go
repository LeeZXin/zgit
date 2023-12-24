package repomd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func GetByPath(ctx context.Context, path string) (Repo, bool, error) {
	var ret Repo
	b, err := xormutil.MustGetXormSession(ctx).Where("path = ?", path).Get(&ret)
	return ret, b, err
}

func UpdateTotalAndGitSize(ctx context.Context, path string, totalSize, gitSize int64) error {
	_, err := xormutil.MustGetXormSession(ctx).Where("path = ?", path).
		Cols("total_size", "git_size").
		Limit(1).
		Update(&Repo{
			TotalSize: totalSize,
			GitSize:   gitSize,
		})
	return err
}

func ListRepo(ctx context.Context, offset int64, limit int, projectId, searchName string) ([]Repo, error) {
	session := xormutil.MustGetXormSession(ctx).Where("project_id = ?", projectId)
	if searchName != "" {
		session.And("name like ?", searchName+"%")
	}
	if offset > 0 {
		session.And("id > ?", offset)
	}
	if limit > 0 {
		session.Limit(limit)
	}
	ret := make([]Repo, 0)
	return ret, session.OrderBy("id asc").Find(&ret)
}

func CountRepo(ctx context.Context, projectId, searchName string) (int64, error) {
	session := xormutil.MustGetXormSession(ctx).Where("project_id = ?", projectId)
	if searchName != "" {
		session.And("name like ?", searchName+"%")
	}
	return session.Count(new(Repo))
}

func UpdateIsEmpty(ctx context.Context, path string, isEmpty bool) error {
	_, err := xormutil.MustGetXormSession(ctx).Where("path = ?", path).
		Cols("is_empty").
		Limit(1).
		Update(&Repo{
			IsEmpty: isEmpty,
		})
	return err
}

func InsertRepo(ctx context.Context, reqDTO InsertRepoReqDTO) (Repo, error) {
	r := Repo{
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

func DeleteRepo(ctx context.Context, path string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("path = ?", path).Limit(1).Delete(new(Repo))
	return rows == 1, err
}
