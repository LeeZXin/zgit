package git

import (
	"context"
	"errors"
	"github.com/LeeZXin/zsf-utils/collections/hashmap"
	"github.com/LeeZXin/zsf-utils/idutil"
	"path/filepath"
	"strconv"
	"strings"
	"zgit/pkg/git/command"
	"zgit/setting"
	"zgit/util"
)

type lsFileLine struct {
	Mode  FileMode
	Sha   string
	Stage int
	Path  string
}

// SameAs checks if two lsFileLines are referring to the same path, sha and mode (ignoring stage)
func (l *lsFileLine) SameAs(other *lsFileLine) bool {
	if l == nil || other == nil {
		return false
	}
	return l.Mode == other.Mode &&
		l.Sha == other.Sha &&
		l.Path == other.Path
}

type unmergedFile struct {
	Path   string
	Based  *lsFileLine
	Ours   *lsFileLine
	Theirs *lsFileLine
}

func (u *unmergedFile) IsConflict() bool {
	if u.Ours != nil && u.Theirs != nil {
		if !u.Ours.SameAs(u.Theirs) {
			return true
		}
		if u.Ours.Mode == SymbolicLinkMode || u.Theirs.Mode == SymbolicLinkMode {
			return true
		}
		if u.Ours.Mode == SubModuleMode || u.Theirs.Mode == SubModuleMode {
			return true
		}
	} else if u.Ours != nil {
		return true
	} else if u.Theirs != nil {
		return true
	}
	return false
}

func findConflictFiles(ctx context.Context, repoPath string, pr DiffCommitsInfo) ([]string, error) {
	tempDir := filepath.Join(setting.TempDir(), "pull-"+idutil.RandomUuid())
	defer util.RemoveAll(tempDir)
	if err := prepare4Merge(ctx, repoPath, tempDir, pr); err != nil {
		return nil, err
	}
	_, err := command.NewCommand("read-tree", "-m", pr.MergeBase, MergeBranch, TrackingBranch).
		Run(ctx, command.WithDir(tempDir))
	if err != nil {
		return nil, err
	}
	unMergedFiles, err := readUnMergedFiles(ctx, tempDir)
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0, len(unMergedFiles))
	for _, file := range unMergedFiles {
		if file.IsConflict() {
			ret = append(ret, file.Path)
		}
	}
	return ret, nil
}

func lsFiles(ctx context.Context, repoPath string, args ...string) ([]lsFileLine, error) {
	readPipe := command.NewCommand("ls-files").AddArgs(args...).RunWithReadPipe(ctx, command.WithDir(repoPath))
	ret := make([]lsFileLine, 0)
	if err := readPipe.RangeStringLines(func(_ int, line string) (bool, error) {
		fields := strings.Fields(line)
		if len(fields) != 4 {
			return false, errors.New("invalid format")
		}
		l := lsFileLine{}
		l.Mode = FileMode(fields[0])
		l.Sha = fields[1]
		stage, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			return false, err
		}
		l.Stage = int(stage)
		l.Path = fields[3]
		ret = append(ret, l)
		return true, nil
	}); err != nil {
		return nil, err
	}
	return ret, nil
}

func lsUnMergedFiles(ctx context.Context, repoPath string) ([]lsFileLine, error) {
	return lsFiles(ctx, repoPath, "-u")
}

func readUnMergedFiles(ctx context.Context, repoPath string) ([]unmergedFile, error) {
	lines, err := lsUnMergedFiles(ctx, repoPath)
	if err != nil {
		return nil, err
	}
	tmp := hashmap.NewHashMap[string, *unmergedFile]()
	for i, line := range lines {
		file, b := tmp.Get(line.Path)
		if !b {
			file = &unmergedFile{
				Path: line.Path,
			}
			tmp.Put(line.Path, file)
		}
		switch line.Stage {
		case 1:
			file.Based = &lines[i]
		case 2:
			file.Ours = &lines[i]
		case 3:
			file.Theirs = &lines[i]
		}
	}
	ret := make([]unmergedFile, 0, tmp.Size())
	tmp.Range(func(_ string, file *unmergedFile) bool {
		ret = append(ret, *file)
		return true
	})
	return ret, nil
}

func UnPackFile(ctx context.Context, repoPath, sha string) (string, error) {
	result, err := command.NewCommand("unpack-file", sha).Run(ctx, command.WithDir(repoPath))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.ReadAsString()), nil
}
