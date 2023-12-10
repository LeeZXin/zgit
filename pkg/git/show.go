package git

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"zgit/pkg/git/command"
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

func GetRefFilesCountsAndSize(ctx context.Context, repoPath, refName string) (int, int64, error) {
	result := command.NewCommand("ls-tree", "--full-tree", "-r", "-l", refName).
		RunWithReadPipe(ctx, command.WithDir(repoPath))
	var (
		fileCount = 0

		fileSize int64 = 0
	)
	if err := result.RangeStringLines(func(_ int, line string) (bool, error) {
		fields := strings.Fields(line)
		if len(fields) == 5 {
			fileCount++
			size, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				fileSize += size
			}
		}
		return true, nil
	}); err != nil {
		return 0, 0, err
	}
	return fileCount, fileSize, nil
}

func CountObjects(ctx context.Context, repoPath string) (int64, float64, error) {
	result, err := command.NewCommand("count-objects").Run(ctx, command.WithDir(repoPath))
	if err != nil {
		return 0, 0, err
	}
	line := result.ReadAsString()
	fields := strings.Fields(strings.TrimSpace(line))
	fmt.Println(line, repoPath)
	if len(fields) != 4 {
		return 0, 0, errors.New("unknown count-objects output")
	}
	objectsCount, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, 0, errors.New("unknown count-objects fields[0]")
	}
	objectsSize, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return 0, 0, errors.New("unknown count-objects fields[2]")
	}
	return objectsCount, objectsSize, nil
}
