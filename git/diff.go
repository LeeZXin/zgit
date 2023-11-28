package git

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"zgit/git/command"
)

var (
	EndBytesFlag = []byte("\000")
)

type CompareInfo struct {
	MergeBase    string
	BaseCommitID string
	HeadCommitID string
	Commits      []*Commit
	NumFiles     int
}

func GetFilesDiffCount(ctx context.Context, repoPath, base, head string, directCompare bool) (int, error) {
	separator := getRefCompareSeparator(directCompare)
	result, err := command.NewCommand("diff", "-z", "--name-only", base+separator+head, "--").Run(ctx, command.WithDir(repoPath))
	if err != nil {
		if strings.Contains(err.Error(), "no merge base") {
			result, err = command.NewCommand("diff", "-z", "--name-only", base, head, "--").Run(ctx, command.WithDir(repoPath))
		}
	}
	if err != nil {
		return 0, err
	}
	return bytes.Count(result.ReadAsBytes(), EndBytesFlag), nil
}

func GetCompareInfoBetween2Ref(ctx context.Context, repoPath, target, head string, directCompare bool) (*CompareInfo, error) {
	if CheckRefIsTag(ctx, repoPath, target) {
		target = TagPrefix + target
	} else if CheckRefIsBranch(ctx, repoPath, target) {
		target = BranchPrefix + target
	} else if CheckRefIsCommit(ctx, repoPath, target) {
		//
	} else {
		return nil, fmt.Errorf("%s is not valid", target)
	}
	if CheckRefIsTag(ctx, repoPath, head) {
		head = TagPrefix + head
	} else if CheckRefIsBranch(ctx, repoPath, head) {
		head = BranchPrefix + head
	} else if CheckRefIsCommit(ctx, repoPath, head) {
		//
	} else {
		return nil, fmt.Errorf("%s is not valid", head)
	}
	var err error
	ret := new(CompareInfo)
	ret.MergeBase, err = MergeBase(ctx, repoPath, target, head)
	if err != nil {
		return nil, err
	}
	ret.BaseCommitID, err = GetRefCommitId(ctx, repoPath, ret.MergeBase)
	return nil, nil
}
