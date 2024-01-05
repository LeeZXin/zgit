package pullrequestmd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func GenPrId() string {
	return idutil.RandomUuid()
}

func GenRid() string {
	return idutil.RandomUuid()
}

func InsertPullRequest(ctx context.Context, reqDTO InsertPullRequestReqDTO) (PullRequest, error) {
	ret := PullRequest{
		PrId:     GenPrId(),
		RepoId:   reqDTO.RepoId,
		Target:   reqDTO.Target,
		Head:     reqDTO.Head,
		PrStatus: reqDTO.PrStatus,
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
			PrStatus: newStatus,
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
			PrStatus:       newStatus,
		})
	return rows == 1, err
}

func GetByPrId(ctx context.Context, prId string) (PullRequest, bool, error) {
	var ret PullRequest
	b, err := xormutil.MustGetXormSession(ctx).Where("pr_id = ?", prId).Get(&ret)
	return ret, b, err
}

func InsertReview(ctx context.Context, reqDTO InsertReviewReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&Review{
		Rid:          GenRid(),
		PrId:         reqDTO.PrId,
		Reviewer:     reqDTO.Reviewer,
		ReviewMsg:    reqDTO.ReviewMsg,
		ReviewStatus: reqDTO.Status,
	})
	return err
}

func UpdateReview(ctx context.Context, reqDTO UpdateReviewReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("rid = ?", reqDTO.Rid).
		Cols("review_status").
		Limit(1).
		Update(&Review{
			ReviewStatus: reqDTO.Status,
		})
	return rows == 1, err
}

func ListReview(ctx context.Context, prId string) ([]Review, error) {
	ret := make([]Review, 0)
	err := xormutil.MustGetXormSession(ctx).Where("pr_id = ?", prId).Find(&ret)
	return ret, err
}

func CountReview(ctx context.Context, prId string, status ReviewStatus) (int, error) {
	ret, err := xormutil.MustGetXormSession(ctx).
		Where("pr_id = ?", prId).
		And("review_status = ?", status.Int()).
		Count(new(Review))
	return int(ret), err
}

func GetReview(ctx context.Context, prId, reviewer string) (Review, bool, error) {
	var ret Review
	b, err := xormutil.MustGetXormSession(ctx).
		Where("pr_id = ?", prId).
		And("reviewer = ?", reviewer).
		Get(&ret)
	return ret, b, err
}
