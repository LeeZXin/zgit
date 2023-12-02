package git

import (
	"context"
	"strings"
	"zgit/git/command"
)

func GetAllTagList(ctx context.Context, repoPath string) ([]string, error) {
	cmd := command.NewCommand("tag")
	pipeResult := cmd.RunWithReadPipe(ctx, command.WithDir(repoPath))
	ret := make([]string, 0)
	err := pipeResult.RangeStringLines(func(_ int, line string) (bool, error) {
		ret = append(ret, strings.TrimSpace(line))
		return true, nil
	})
	return ret, err
}

func CheckRefIsTag(ctx context.Context, repoPath string, tag string) bool {
	if !strings.HasPrefix(tag, TagPrefix) {
		tag = TagPrefix + tag
	}
	return CatFileExists(ctx, repoPath, tag) == nil
}
