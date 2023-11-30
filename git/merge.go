package git

import (
	"context"
	"strings"
	"zgit/git/command"
)

func MergeBase(ctx context.Context, repoPath, target, head string) (string, error) {
	result, err := command.NewCommand("merge-base", "--", target, head).Run(ctx, command.WithDir(repoPath))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.ReadAsString()), err
}

func Merge(ctx context.Context, repoPath, target, head string) error {
	return nil
}
