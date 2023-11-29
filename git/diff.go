package git

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"zgit/git/command"
)

type DiffFileType int

const (
	ModifiedFileType DiffFileType = iota
	CreatedFileType
	DeletedFileType
	RenamedFileType
	CopiedFileType
)

const (
	InsertLinePrefix = "+"
	DeleteLinePrefix = "-"
	NormalLinePrefix = " "
)

func (t DiffFileType) String() string {
	switch t {
	case CreatedFileType:
		return "created"
	case DeletedFileType:
		return "deleted"
	case RenamedFileType:
		return "renamed"
	case CopiedFileType:
		return "copied"
	case ModifiedFileType:
		return "modified"
	default:
		return "unknown"
	}
}

var (
	EndBytesFlag = []byte("\000")
)

const (
	DefaultFileMode = "100644"
)

type CompareInfo struct {
	TargetCommitId string    `json:"targetCommitId"`
	HeadCommitId   string    `json:"headCommitId"`
	Commits        []*Commit `json:"commits"`
	NumFiles       int       `json:"numFiles"`
}

type DiffNumsStat struct {
	Path       string
	TotalNums  int
	InsertNums int
	DeleteNums int
}

type DiffDetail struct {
	OldMode     string
	Mode        string
	IsSubModule bool
	FileType    DiffFileType
	IsBinary    bool
}

func NewDiffDetail() *DiffDetail {
	return &DiffDetail{
		OldMode:  DefaultFileMode,
		Mode:     DefaultFileMode,
		FileType: ModifiedFileType,
	}
}

type DiffLine struct {
	Index   int    `json:"index"`
	LeftNo  int    `json:"leftNo"`
	Prefix  string `json:"prefix"`
	RightNo int    `json:"rightNo"`
	Text    string `json:"text"`
}

func GetFilesDiffCount(ctx context.Context, repoPath, target, head string) (int, error) {
	result, err := command.NewCommand("diff", "-z", "--name-only", target+".."+head, "--").Run(ctx, command.WithDir(repoPath))
	if err != nil {
		if strings.Contains(err.Error(), "no merge base") {
			result, err = command.NewCommand("diff", "-z", "--name-only", target, head, "--").Run(ctx, command.WithDir(repoPath))
		}
	}
	if err != nil {
		return 0, err
	}
	return bytes.Count(result.ReadAsBytes(), EndBytesFlag), nil
}

func GetCompareInfoBetween2Ref(ctx context.Context, repoPath, target, head string) (*CompareInfo, error) {
	var err error
	if CheckRefIsTag(ctx, repoPath, target) {
		target = TagPrefix + target
	} else if CheckRefIsBranch(ctx, repoPath, target) {
		target = BranchPrefix + target
	} else if CheckRefIsCommit(ctx, repoPath, target) {
		target, err = GetFullShaCommitId(ctx, repoPath, target)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("%s is not valid", target)
	}
	if CheckRefIsTag(ctx, repoPath, head) {
		head = TagPrefix + head
	} else if CheckRefIsBranch(ctx, repoPath, head) {
		head = BranchPrefix + head
	} else if CheckRefIsCommit(ctx, repoPath, head) {
		head, err = GetFullShaCommitId(ctx, repoPath, head)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("%s is not valid", head)
	}
	ret := new(CompareInfo)
	ret.TargetCommitId, err = GetRefCommitId(ctx, repoPath, target)
	if err != nil {
		return nil, err
	}
	ret.HeadCommitId, err = GetRefCommitId(ctx, repoPath, head)
	if err != nil {
		return nil, err
	}
	ret.Commits, err = GetGitLogCommitList(ctx, repoPath, ret.HeadCommitId, ret.TargetCommitId)
	if err != nil {
		return nil, err
	}
	ret.NumFiles, err = GetFilesDiffCount(ctx, repoPath, ret.HeadCommitId, ret.TargetCommitId)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func GetDiffNumsStat(ctx context.Context, repoPath, target, head string) ([]DiffNumsStat, error) {
	pipeResult := command.NewCommand("diff", "--numstat", target+".."+head, "--").RunWithReadPipe(ctx, command.WithDir(repoPath))
	defer pipeResult.ClosePipe()
	reader := bufio.NewReader(pipeResult.Reader())
	ret := make([]DiffNumsStat, 0)
	for {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if isPrefix {
			continue
		}
		lineStr := strings.TrimSpace(string(line))
		fields := strings.Fields(lineStr)
		if len(fields) != 3 {
			continue
		}
		insertNums, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("parseInt err: %v", insertNums)
		}
		deleteNums, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("parseInt err: %v", insertNums)
		}
		ret = append(ret, DiffNumsStat{
			Path:       fields[2],
			InsertNums: insertNums,
			DeleteNums: deleteNums,
			TotalNums:  insertNums + deleteNums,
		})
	}
	return ret, nil
}

func GetDiffShortStat(ctx context.Context, repoPath, target, head string) (int, int, int, error) {
	result, err := command.NewCommand("diff", "--shortstat", target+".."+head, "--").Run(ctx, command.WithDir(repoPath))
	if err != nil {
		return 0, 0, 0, err
	}
	var (
		fileChangeNums, insertNums, deleteNums int
	)
	line := strings.TrimSpace(result.ReadAsString())
	lineSplit := strings.Split(line, ",")
	for _, item := range lineSplit {
		fields := strings.Fields(item)
		if strings.Contains(item, "files changed") {
			fileChangeNums, err = strconv.Atoi(fields[0])
			if err != nil {
				return 0, 0, 0, fmt.Errorf("parseInt err:%v", err)
			}
		} else if strings.Contains(item, "insertions") {
			insertNums, err = strconv.Atoi(fields[0])
			if err != nil {
				return 0, 0, 0, fmt.Errorf("parseInt err:%v", err)
			}
		} else if strings.Contains(item, "deletions") {
			deleteNums, err = strconv.Atoi(fields[0])
			if err != nil {
				return 0, 0, 0, fmt.Errorf("parseInt err:%v", err)
			}
		}
	}
	return fileChangeNums, insertNums, deleteNums, nil
}

func GetDiffFileContent(ctx context.Context, repoPath, target, head, filePath string) (*DiffDetail, error) {
	pipeResult := command.NewCommand("diff", "--src-prefix=a/", "--dst-prefix=b/", target+".."+head, "--", filePath).RunWithReadPipe(ctx, command.WithDir(repoPath))
	defer pipeResult.ClosePipe()
	reader := bufio.NewReaderSize(pipeResult.Reader(), 4096)
	c := NewDiffDetail()
	for {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if isPrefix {
			continue
		}
		lineStr := strings.TrimSpace(string(line))
		if strings.HasPrefix(lineStr, "diff --git") {
			// nothing
		} else if strings.HasPrefix(lineStr, "index") {
			if strings.HasSuffix(lineStr, "160000") {
				c.IsSubModule = true
			}
		} else if strings.HasPrefix(lineStr, "---") {
			// nothing
		} else if strings.HasPrefix(lineStr, "+++") {
			// parse hunks
		} else if strings.HasPrefix(lineStr, "new mode") {
			c.Mode = strings.TrimSpace(strings.TrimPrefix(lineStr, "new mode"))
			if strings.HasSuffix(lineStr, "160000") {
				c.IsSubModule = true
			}
		} else if strings.HasPrefix(lineStr, "old mode") {
			c.OldMode = strings.TrimSpace(strings.TrimPrefix(lineStr, "old mode"))
			if strings.HasSuffix(lineStr, "160000") {
				c.IsSubModule = true
			}
		} else if strings.HasPrefix(lineStr, "new file mode") {
			c.FileType = CreatedFileType
			c.Mode = strings.TrimSpace(strings.TrimPrefix(lineStr, "new mode"))
			if strings.HasSuffix(lineStr, "160000") {
				c.IsSubModule = true
			}
		} else if strings.HasPrefix(lineStr, "rename from") {
			c.FileType = RenamedFileType
		} else if strings.HasPrefix(lineStr, "rename to") {
			c.FileType = RenamedFileType
		} else if strings.HasPrefix(lineStr, "copy from") {
			c.FileType = CopiedFileType
		} else if strings.HasPrefix(lineStr, "copy to") {
			c.FileType = CopiedFileType
		} else if strings.HasPrefix(lineStr, "deleted") {
			c.FileType = DeletedFileType
		} else if strings.HasPrefix(lineStr, "similarity index") {
			c.FileType = RenamedFileType
		} else if strings.HasPrefix(lineStr, "Binary") {
			c.IsBinary = true
		}
	}
	return c, nil
}

func parseHunks(ctx context.Context, reader *bufio.Reader) ([]DiffLine, error) {
	insertionNums := 0
	deletionNums := 0
	totalLinesNums := 0
	ret := make([]DiffLine, 0)
	var (
		index, leftNo, rightNo int
	)
	for {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if isPrefix {
			continue
		}
		totalLinesNums++
		lineStr := strings.TrimSpace(string(line))
		switch lineStr[0] {
		case '@':
			leftNo, _, rightNo, _, err = parseHunkString(lineStr)
			if err != nil {
				return nil, err
			}
		case '+':
			if len(ret) > 0 && ret[len(ret)-1].Prefix == DeleteLinePrefix {
				ret = append(ret, DiffLine{
					Index:   index,
					LeftNo:  leftNo,
					Prefix:  InsertLinePrefix,
					RightNo: rightNo,
					Text:    lineStr,
				})
			} else {
				rightNo++
				ret = append(ret, DiffLine{
					Index:   index,
					LeftNo:  leftNo,
					Prefix:  InsertLinePrefix,
					RightNo: rightNo,
					Text:    lineStr,
				})
			}

			index++
			insertionNums++
		case '-':
			leftNo++
			ret = append(ret, DiffLine{
				Index:   index,
				LeftNo:  leftNo,
				Prefix:  DeleteLinePrefix,
				RightNo: rightNo,
				Text:    lineStr,
			})
			index++
			deletionNums++
		default:
			leftNo++
			rightNo++
			ret = append(ret, DiffLine{
				Index:   index,
				LeftNo:  leftNo,
				Prefix:  NormalLinePrefix,
				RightNo: rightNo,
				Text:    lineStr,
			})
			index++
		}
	}
	return ret, nil
}

func parseHunkString(line string) (int, int, int, int, error) {
	fields := strings.Fields(line)
	if len(fields) < 4 || fields[0] != "@@" || fields[3] != "@@" {
		return 0, 0, 0, 0, errors.New("invalid @@ format")
	}
	o := strings.Split(fields[1][1:], ",")
	if len(o) < 2 {
		return 0, 0, 0, 0, errors.New("invalid @@ format")
	}
	o1, err := strconv.Atoi(o[0])
	if err != nil {
		return 0, 0, 0, 0, err
	}
	o2, err := strconv.Atoi(o[1])
	if err != nil {
		return 0, 0, 0, 0, err
	}
	n := strings.Split(fields[2][1:], ",")
	if len(o) < 2 {
		return 0, 0, 0, 0, errors.New("invalid @@ format")
	}
	n1, err := strconv.Atoi(n[0])
	if err != nil {
		return 0, 0, 0, 0, err
	}
	n2, err := strconv.Atoi(n[1])
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return o1, o2, n1, n2, nil
}
