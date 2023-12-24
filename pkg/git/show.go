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

func GetRefFilesCountAndSize(ctx context.Context, repoPath, refName string) (int, int64, error) {
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

type LsTreeRet struct {
	Mode FileMode
	Path string
	Size int64
	Blob string
}

func LsTreeWithoutRecurse(ctx context.Context, repoPath, refName, dir string, offset, limit int) ([]LsTreeRet, error) {
	cmd := command.NewCommand("ls-tree", "--full-tree", "-l", refName)
	if dir != "" {
		cmd.AddArgs("--", dir+"/")
	}
	result := cmd.RunWithReadPipe(ctx, command.WithDir(repoPath))
	if offset < 0 {
		offset = 0
	}
	var endNum int
	if limit < 0 {
		endNum = -1
	} else {
		endNum = offset + limit
	}
	ret := make([]LsTreeRet, 0)
	if err := result.RangeStringLines(func(n int, line string) (bool, error) {
		if n < offset {
			return true, nil
		}
		if endNum > 0 && n >= endNum {
			return false, nil
		}
		fields := strings.Fields(line)
		if len(fields) == 5 {
			size, _ := strconv.ParseInt(fields[3], 10, 64)
			ret = append(ret, LsTreeRet{
				Mode: FileMode(fields[0]),
				Path: fields[4],
				Size: size,
				Blob: fields[2],
			})
		}
		return true, nil
	}); err != nil {
		return nil, err
	}
	return ret, nil
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

type FileCommit struct {
	LsTreeRet
	*Commit
}

func LsTreeCommit(ctx context.Context, repoPath, refName string, dir string, offset, limit int) ([]FileCommit, error) {
	lsRet, err := LsTreeWithoutRecurse(ctx, repoPath, refName, dir, offset, limit)
	if err != nil {
		return nil, err
	}
	commits := make([]FileCommit, 0, len(lsRet))
	for _, ret := range lsRet {
		commit, err := GetFileLastCommit(ctx, repoPath, refName, ret.Path)
		if err != nil {
			return nil, err
		}
		commits = append(commits, FileCommit{
			LsTreeRet: ret,
			Commit:    commit,
		})
	}
	return commits, nil
}

func GetFileContentByBlob(ctx context.Context, repoPath, blob string) (string, error) {
	result, err := command.NewCommand("show", blob).Run(ctx, command.WithDir(repoPath))
	if err != nil {
		return "", err
	}
	return result.ReadAsString(), nil
}

func GetFileContentByRef(ctx context.Context, repoPath, refName, filePath string) (FileMode, string, bool, error) {
	result, err := command.NewCommand("ls-tree", "--full-tree", "-l", refName, "--", filePath).Run(ctx, command.WithDir(repoPath))
	if err != nil {
		return "", "", false, err
	}
	ret := result.ReadAsString()
	if ret == "" {
		return "", "", false, nil
	}
	fields := strings.Fields(strings.TrimSpace(ret))
	if len(fields) < 3 {
		return "", "", false, errors.New("unknown format")
	}
	blob := fields[2]
	mode := fields[0]
	if mode == RegularFileMode.String() {
		content, err := GetFileContentByBlob(ctx, repoPath, blob)
		if err != nil {
			return "", "", false, err
		}
		return RegularFileMode, content, true, nil
	}
	return FileMode(mode), "", true, nil
}
