package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/LeeZXin/zsf-utils/idutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"zgit/pkg/git/command"
	"zgit/setting"
	"zgit/util"
)

const (
	MergeBranch    = "base"
	TrackingBranch = "tracking"
)

var (
	escapedSymbols = regexp.MustCompile(`([*[?! \\])`)
)

type PreparePullRequestInfo struct {
	OriginHead    string           `json:"originHead"`
	OriginTarget  string           `json:"originTarget"`
	Target        string           `json:"target"`
	Head          string           `json:"head"`
	TargetCommit  Commit           `json:"targetCommit"`
	HeadCommit    Commit           `json:"headCommit"`
	Commits       []Commit         `json:"commits"`
	NumFiles      int              `json:"numFiles"`
	DiffNumsStats DiffNumsStatInfo `json:"diffNumsStats"`
}

type MergeRepoOpts struct {
	Message string
}

func PreparePullRequest(ctx context.Context, repoPath, target, head string) (PreparePullRequestInfo, error) {
	pr := PreparePullRequestInfo{}
	pr.OriginTarget, pr.OriginHead = target, head
	if !strings.HasPrefix(head, BranchPrefix) {
		head = BranchPrefix + head
	}
	if !CheckRefIsBranch(ctx, repoPath, head) {
		return PreparePullRequestInfo{}, fmt.Errorf("%s is not branch", head)
	}
	commitId, err := GetRefCommitId(ctx, repoPath, head)
	if err != nil {
		return PreparePullRequestInfo{}, err
	}
	pr.HeadCommit, err = GetCommitByCommitId(ctx, repoPath, commitId)
	if err != nil {
		return PreparePullRequestInfo{}, err
	}
	if CheckRefIsTag(ctx, repoPath, target) {
		if !strings.HasPrefix(target, TagPrefix) {
			target = TagPrefix + target
		}
		pr.TargetCommit, err = GetCommitByTag(ctx, repoPath, target)
		if err != nil {
			return PreparePullRequestInfo{}, err
		}
	} else if CheckRefIsBranch(ctx, repoPath, target) {
		if !strings.HasPrefix(target, BranchPrefix) {
			target = BranchPrefix + target
		}
		commitId, err = GetRefCommitId(ctx, repoPath, target)
		if err != nil {
			return PreparePullRequestInfo{}, err
		}
		pr.TargetCommit, err = GetCommitByCommitId(ctx, repoPath, commitId)
		if err != nil {
			return PreparePullRequestInfo{}, err
		}
	} else if CheckRefIsCommit(ctx, repoPath, target) {
		commitId, err = GetRefCommitId(ctx, repoPath, target)
		if err != nil {
			return PreparePullRequestInfo{}, err
		}
		pr.TargetCommit, err = GetCommitByCommitId(ctx, repoPath, commitId)
		if err != nil {
			return PreparePullRequestInfo{}, err
		}
	} else {
		return PreparePullRequestInfo{}, fmt.Errorf("%s unsupported type", target)
	}
	// 这里要反过来 git log 查看target的提交记录 不是head的提交记录
	pr.Commits, err = GetGitLogCommitList(ctx, repoPath, pr.HeadCommit.Id, pr.TargetCommit.Id)
	if err != nil {
		return PreparePullRequestInfo{}, err
	}
	pr.NumFiles, err = GetFilesDiffCount(ctx, repoPath, pr.TargetCommit.Id, pr.HeadCommit.Id)
	if err != nil {
		return PreparePullRequestInfo{}, err
	}
	pr.DiffNumsStats, err = GenDiffNumsStat(ctx, repoPath, pr.TargetCommit.Id, pr.HeadCommit.Id)
	if err != nil {
		return PreparePullRequestInfo{}, err
	}
	pr.Head, pr.Target = head, target
	return pr, nil
}

func Merge(ctx context.Context, repoPath, target, head string, opts MergeRepoOpts) error {
	prInfo, err := PreparePullRequest(ctx, repoPath, target, head)
	if err != nil {
		return err
	}
	return doMerge(ctx, repoPath, prInfo, opts)
}

func doMerge(ctx context.Context, repoPath string, pr PreparePullRequestInfo, opts MergeRepoOpts) error {
	if len(pr.Commits) == 0 {
		return errors.New("nothing to commit")
	}
	tempDir := filepath.Join(setting.TempDir(), "merge-"+idutil.RandomUuid())
	defer util.RemoveAll(tempDir)
	var err error
	if err = initEmptyRepository(ctx, tempDir, false); err != nil {
		return err
	}
	if _, err = command.NewCommand("remote", "add", "-t", pr.OriginHead, "-m", pr.OriginHead, "origin", repoPath).
		Run(ctx, command.WithDir(tempDir)); err != nil {
		return errors.New("add remote failed")
	}
	fetchArgs := make([]string, 0)
	fetchArgs = append(fetchArgs, "--no-tags")
	if CheckGitVersionAtLeast("2.25.0") == nil {
		fetchArgs = append(fetchArgs, "--no-write-commit-graph")
	}
	if _, err = command.NewCommand("fetch", "origin", pr.OriginHead+":"+MergeBranch, pr.OriginHead+":original_"+pr.OriginHead).AddArgs(fetchArgs...).
		Run(ctx, command.WithDir(tempDir)); err != nil {
		return err
	}
	if err = SetDefaultBranch(ctx, tempDir, MergeBranch); err != nil {
		return err
	}
	if _, err = command.NewCommand("fetch", "origin", pr.TargetCommit.Id+":"+TrackingBranch).AddArgs(fetchArgs...).
		Run(ctx, command.WithDir(tempDir)); err != nil {
		return err
	}
	infoPath := filepath.Join(tempDir, ".git", "info")
	if err = os.MkdirAll(infoPath, 0o700); err != nil {
		return fmt.Errorf("unable to create .git/info in tmpBasePath: %w", err)
	}
	sparseCheckout, err := os.OpenFile(filepath.Join(infoPath, "sparse-checkout"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("unable to write .git/info/sparse-checkout file in tmpBasePath: %w", err)
	}
	defer sparseCheckout.Close()
	trees, err := getDiffTreeForMerge(ctx, tempDir, TrackingBranch, MergeBranch)
	if err != nil {
		return fmt.Errorf("unable to get diff tree in tmpBasePath: %w", err)
	}
	for _, tree := range trees {
		if _, err = sparseCheckout.WriteString(tree); err != nil {
			return fmt.Errorf("unable to write to sparseCheckout in tmpBasePath: %w", err)
		}
	}
	if err = SetLocalConfig(ctx, tempDir, "filter.lfs.process", ""); err != nil {
		return err
	}
	if err = SetLocalConfig(ctx, tempDir, "filter.lfs.required", "false"); err != nil {
		return err
	}
	if err = SetLocalConfig(ctx, tempDir, "filter.lfs.clean", ""); err != nil {
		return err
	}
	if err = SetLocalConfig(ctx, tempDir, "filter.lfs.smudge", ""); err != nil {
		return err
	}
	if err = SetLocalConfig(ctx, tempDir, "core.sparseCheckout", "true"); err != nil {
		return err
	}
	if _, err = command.NewCommand("read-tree", "HEAD").Run(ctx, command.WithDir(tempDir)); err != nil {
		return err
	}
	if _, err = command.NewCommand("merge", "--no-ff", "--no-commit", TrackingBranch).
		Run(ctx, command.WithDir(tempDir)); err != nil {
		if _, statErr := os.Stat(filepath.Join(tempDir, ".git", "MERGE_HEAD")); statErr == nil {
			return &ErrMergeConflict{
				err: err,
			}
		} else if strings.Contains(err.Error(), "refusing to merge unrelated histories") {
			return &ErrMergeUnrelatedHistories{
				err: err,
			}
		}
		return fmt.Errorf("git merge err: %v", err)
	}
	mergeCmd := command.NewCommand("commit", "--no-gpg-sign", "-m", opts.Message)
	if _, err = mergeCmd.Run(ctx, command.WithDir(tempDir)); err != nil {
		return err
	}
	if _, err = command.NewCommand("push", "origin", MergeBranch+":"+pr.Head).
		Run(ctx, command.WithDir(tempDir)); err != nil {
		if strings.Contains(err.Error(), "non-fast-forward") {
			return &ErrPushOutOfDate{
				err: err,
			}
		} else if strings.Contains(err.Error(), "! [remote rejected]") {
			return &ErrPushRejected{
				err: err,
			}
		}
		return fmt.Errorf("git push: %v", err)
	}
	return nil
}

func getDiffTreeForMerge(ctx context.Context, repoPath, target, head string) ([]string, error) {
	diffTreeResult, err := command.NewCommand("diff-tree", "--no-commit-id", "--name-only", "-r", "-r", "-z", "--root", target, head).
		Run(ctx, command.WithDir(repoPath))
	if err != nil {
		return nil, fmt.Errorf("unable to diff tree in tmpBasePath: %w", err)
	}
	treeResult := bytes.Split(diffTreeResult.ReadAsBytes(), []byte{'\x00'})
	ret := make([]string, 0)
	for _, r := range treeResult {
		line := strings.TrimSpace(string(r))
		if len(line) > 0 {
			ret = append(ret, fmt.Sprintf("/%s\n", escapedSymbols.ReplaceAllString(line, `\$1`)))
		}
	}
	return ret, nil
}
