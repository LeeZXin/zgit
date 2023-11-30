package git

import (
	"context"
	"strings"
	"zgit/git/command"
)

func GetAllBranchList(ctx context.Context, repoPath string) ([]string, error) {
	cmd := command.NewCommand("for-each-ref", "--format=%(objectname) %(refname)", BranchPrefix, "--sort=-committerdate")
	pipeResult := cmd.RunWithReadPipe(ctx, command.WithDir(repoPath))
	defer pipeResult.ClosePipe()
	ret := make([]string, 0)
	err := pipeResult.RangeStringLines(func(_ int, line string) error {
		split := strings.Split(strings.TrimSpace(line), " ")
		if len(split) == 2 {
			ret = append(ret, strings.TrimPrefix(split[1], BranchPrefix))
		}
		return nil
	})
	return ret, err
}

func NewBranch(ctx context.Context, repoPath string, name string) error {
	_, err := command.NewCommand("branch", name).Run(ctx, command.WithDir(repoPath))
	return err
}

func CheckRefIsBranch(ctx context.Context, repoPath string, branch string) bool {
	return CatFileExists(ctx, repoPath, BranchPrefix+branch) == nil
}

func CheckCommitIfInBranch(ctx context.Context, repoPath, commitId, branch string) (bool, error) {
	result, err := command.NewCommand("branch", "--contains", commitId, branch).Run(ctx, command.WithDir(repoPath))
	if err != nil {
		return false, err
	}
	return len(strings.TrimSpace(result.ReadAsString())) > 0, nil
}
