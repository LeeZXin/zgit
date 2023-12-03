package git

import (
	"context"
	"strings"
	"zgit/git/command"
)

func GetAllBranchList(ctx context.Context, repoPath string) ([]string, error) {
	cmd := command.NewCommand("for-each-ref", "--format=%(objectname) %(refname)", BranchPrefix, "--sort=-committerdate")
	pipeResult := cmd.RunWithReadPipe(ctx, command.WithDir(repoPath))
	ret := make([]string, 0)
	err := pipeResult.RangeStringLines(func(_ int, line string) (bool, error) {
		split := strings.Split(strings.TrimSpace(line), " ")
		if len(split) == 2 {
			ret = append(ret, strings.TrimPrefix(split[1], BranchPrefix))
		}
		return true, nil
	})
	return ret, err
}

func CheckRefIsBranch(ctx context.Context, repoPath string, branch string) bool {
	if !strings.HasPrefix(branch, BranchPrefix) {
		branch = BranchPrefix + branch
	}
	return CatFileExists(ctx, repoPath, branch) == nil
}

func CheckCommitIfInBranch(ctx context.Context, repoPath, commitId, branch string) (bool, error) {
	result, err := command.NewCommand("branch", "--contains", commitId, branch).Run(ctx, command.WithDir(repoPath))
	if err != nil {
		return false, err
	}
	return len(strings.TrimSpace(result.ReadAsString())) > 0, nil
}
