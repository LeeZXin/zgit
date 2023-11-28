package git

import (
	"context"
	"strings"
	"zgit/git/command"
)

func GetAllTagList(ctx context.Context, repoPath string) ([]string, error) {
	cmd := command.NewCommand("tag")
	pipeResult := cmd.RunWithReadPipe(ctx, command.WithDir(repoPath))
	defer pipeResult.ClosePipe()
	ret := make([]string, 0)
	err := pipeResult.RangeStringLines(func(_ int, line string) error {
		ret = append(ret, strings.TrimSpace(line))
		return nil
	})
	return ret, err
}

func CheckRefIsTag(ctx context.Context, repoPath string, tag string) bool {
	return CatFileExists(ctx, repoPath, TagPrefix+tag) == nil
}
