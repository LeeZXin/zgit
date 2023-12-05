package git

import (
	"context"
	"fmt"
	"zgit/git/command"
)

func ShowFileTextContentByCommitId(ctx context.Context, repoPath, commitId, filePath string, startLine, limit int) ([]string, error) {
	pipeResult := command.NewCommand("show", fmt.Sprintf("%s:%s", commitId, filePath), "--text").
		RunWithReadPipe(ctx, command.WithDir(repoPath))
	ret := make([]string, 0)
	endLine := startLine + limit
	if err := pipeResult.RangeStringLines(func(i int, line string) (bool, error) {
		if i < startLine {
			return true, nil
		}
		if limit < 0 || (i >= startLine && i < endLine) {
			ret = append(ret, line)
			return true, nil
		}
		return false, nil
	}); err != nil {
		return nil, err
	}
	return ret, nil
}
