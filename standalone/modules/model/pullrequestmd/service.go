package pullrequestmd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func GenPrId() string {
	return idutil.RandomUuid()
}

func InsertPullRequest(ctx context.Context, reqDTO InsertPullRequestReqDTO) (PullRequest, error) {
	ret := PullRequest{
		PrId:     GenPrId(),
		RepoId:   reqDTO.RepoId,
		Target:   reqDTO.Target,
		Head:     reqDTO.Head,
		PrStatus: reqDTO.PrStatus.Int(),
		CreateBy: reqDTO.CreateBy,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func UpdatePrStatus(ctx context.Context, prId string, oldStatus, newStatus PrStatus) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("pr_id = ?", prId).
		And("pr_status = ?", oldStatus.Int()).
		Cols("pr_status").
		Update(&PullRequest{
			PrStatus: newStatus.Int(),
		})
	return rows == 1, err
}

func UpdatePrStatusAndCommitId(ctx context.Context, prId string, oldStatus, newStatus PrStatus, targetCommitId, headCommitId string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("pr_id = ?", prId).
		And("pr_status = ?", oldStatus.Int()).
		Cols("pr_status", "target_commit_id", "head_commit_id").
		Update(&PullRequest{
			TargetCommitId: targetCommitId,
			HeadCommitId:   headCommitId,
			PrStatus:       newStatus.Int(),
		})
	return rows == 1, err
}

func GetByPrId(ctx context.Context, prId string) (PullRequest, bool, error) {
	var ret PullRequest
	b, err := xormutil.MustGetXormSession(ctx).Where("pr_id = ?", prId).Get(&ret)
	return ret, b, err
}
