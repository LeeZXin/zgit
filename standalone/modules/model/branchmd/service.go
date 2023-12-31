package branchmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func InsertProtectedBranch(ctx context.Context, repoPath, branch string) (ProtectedBranch, error) {
	ret := ProtectedBranch{
		Branch: branch,
		RepoId: repoPath,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func DeleteProtectedBranch(ctx context.Context, branch ProtectedBranch) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", branch.RepoId).
		And("branch = ?", branch.Branch).
		Limit(1).
		Delete(new(ProtectedBranch))
	return rows == 1, err
}

func GetProtectedBranch(ctx context.Context, repoPath, branch string) (ProtectedBranch, bool, error) {
	ret := ProtectedBranch{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", repoPath).
		And("branch = ?", branch).
		Limit(1).
		Get(&ret)
	return ret, b, err
}

func ListProtectedBranch(ctx context.Context, reqDTO ListProtectedBranchReqDTO) ([]ProtectedBranch, error) {
	session := xormutil.MustGetXormSession(ctx).Where("repo_id = ?", reqDTO.RepoId)
	if reqDTO.SearchName != "" {
		session.And("branch like ?", reqDTO.SearchName+"%")
	}
	if reqDTO.Offset > 0 {
		session.And("id > ?", reqDTO.Offset)
	}
	if reqDTO.Limit > 0 {
		session.Limit(reqDTO.Limit)
	}
	ret := make([]ProtectedBranch, 0)
	err := session.Find(&ret)
	return ret, err
}

func CountProtectedBranch(ctx context.Context, reqDTO ListProtectedBranchReqDTO) (int64, error) {
	session := xormutil.MustGetXormSession(ctx).Where("repo_id = ?", reqDTO.RepoId)
	if reqDTO.SearchName != "" {
		session.And("branch like ?", reqDTO.SearchName+"%")
	}
	return session.Count(new(ProtectedBranch))
}
