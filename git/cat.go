package git

import (
	"bufio"
	"context"
	"fmt"
	"strconv"
	"strings"
	"zgit/git/command"
)

func BatchCheck(ctx context.Context, repoPath string, name string) (string, string, int64, error) {
	cmd := command.NewCommand("cat-file", "--batch-check")
	pipe := cmd.RunWithStdinPipe(ctx, command.WithDir(repoPath))
	defer pipe.ClosePipe()
	_, err := pipe.Writer().Write([]byte(name + "\n"))
	if err != nil {
		return "", "", 0, err
	}
	reader := bufio.NewReader(pipe.Reader())
	for {
		line, isPrefix, err := reader.ReadLine()
		if isPrefix {
			continue
		}
		if err != nil {
			return "", "", 0, fmt.Errorf("readline err: %v", err)
		}
		fields := strings.Fields(string(line))
		if len(fields) != 3 {
			return "", "", 0, fmt.Errorf("%s does not exists", name)
		}
		size, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			return "", "", 0, fmt.Errorf("parse size error: %v", err)
		}
		return fields[0], fields[1], size, nil
	}
}

func GetRefCommitId(ctx context.Context, repoPath string, name string) (string, error) {
	commitId, _, _, err := BatchCheck(ctx, repoPath, name)
	return commitId, err
}
