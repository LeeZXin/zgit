package repomd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
)

func GenRepoId() string {
	return idutil.RandomUuid()
}

func GetByPath(ctx context.Context, path string) (Repo, bool, error) {
	var ret Repo
	b, err := mysqlstore.GetXormSession(ctx).Where("path = ?", path).Get(&ret)
	return ret, b, err
}

func GetByRepoId(ctx context.Context, repoId string) (Repo, bool, error) {
	var ret Repo
	b, err := mysqlstore.GetXormSession(ctx).Where("repo_id = ?", repoId).Get(&ret)
	return ret, b, err
}

func UpdateTotalAndGitSize(ctx context.Context, repoId string, totalSize, gitSize int64) error {
	_, err := mysqlstore.GetXormSession(ctx).Where("repo_id = ?", repoId).
		Cols("total_size", "git_size").
		Limit(1).
		Update(&Repo{
			TotalSize: totalSize,
			GitSize:   gitSize,
		})
	return err
}

func InsertRepo(ctx context.Context, reqDTO InsertRepoReqDTO) (Repo, error) {
	r := Repo{
		RepoId:        GenRepoId(),
		Name:          reqDTO.Name,
		Path:          reqDTO.Path,
		UserId:        reqDTO.UserId,
		NodeId:        reqDTO.NodeId,
		CorpId:        reqDTO.CorpId,
		ProjectId:     reqDTO.ProjectId,
		RepoDesc:      reqDTO.RepoDesc,
		DefaultBranch: reqDTO.DefaultBranch,
		RepoType:      int(reqDTO.RepoType),
		IsEmpty:       reqDTO.IsEmpty,
		TotalSize:     reqDTO.TotalSize,
		GitSize:       reqDTO.GitSize,
		LfsSize:       reqDTO.LfsSize,
	}
	_, err := mysqlstore.GetXormSession(ctx).Insert(&r)
	return r, err
}

func DeleteRepo(ctx context.Context, repoId string) (bool, error) {
	rows, err := mysqlstore.GetXormSession(ctx).Where("repo_id = ?", repoId).Limit(1).Delete(new(Repo))
	return rows == 1, err
}
