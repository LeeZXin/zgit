package git

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"zgit/pkg/git/command"
)

func CatFileBatchCheck(ctx context.Context, repoPath string, name string) (string, string, int64, error) {
	cmd := command.NewCommand("cat-file", "--batch-check")
	pipe := cmd.RunWithStdinPipe(ctx, command.WithDir(repoPath))
	_, err := pipe.Writer().Write([]byte(name + "\n"))
	if err != nil {
		return "", "", 0, err
	}
	var (
		ref, typ string
		size     int64
		e        error
	)
	if err = pipe.RangeStringLines(func(_ int, line string) (bool, error) {
		ref, typ, size, e = readBatchLine(line)
		return false, nil
	}); err != nil {
		return "", "", 0, err
	}
	return ref, typ, size, e
}

func readBatchLine(line string) (string, string, int64, error) {
	fields := strings.Fields(line)
	if len(fields) == 2 {
		return fields[0], "", 0, errors.New("ref is missing")
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

func CheckExists(ctx context.Context, repoPath string, name string) bool {
	_, err := command.NewCommand("cat-file", "-e", name).Run(ctx, command.WithDir(repoPath))
	return err == nil
}
