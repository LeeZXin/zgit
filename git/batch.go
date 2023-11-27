package git

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"zgit/git/command"
)

func CatFileBatchCheck(ctx context.Context, repoPath string, name string) (string, string, int64, error) {
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
		return readBatchLine(string(line))
	}
}

func readBatchLine(line string) (string, string, int64, error) {
	fields := strings.Fields(line)
	if len(fields) == 2 {
		return fields[0], fields[1], 0, nil
	}
	if len(fields) != 3 {
		return "", "", 0, errors.New("format error")
	}
	size, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return "", "", 0, fmt.Errorf("parse size error: %v", err)
	}
	return fields[0], fields[1], size, nil
}

func CatFileBatch(ctx context.Context, repoPath string, name string, readFn func(io.Reader, command.PipeResultCloser) error) error {
	if readFn == nil {
		return errors.New("readFn is nil")
	}
	cmd := command.NewCommand("cat-file", "--batch")
	pipe := cmd.RunWithStdinPipe(ctx, command.WithDir(repoPath))
	defer pipe.ClosePipe()
	_, err := pipe.Writer().Write([]byte(name + "\n"))
	if err != nil {
		return err
	}
	return readFn(pipe.Reader(), pipe)
}
