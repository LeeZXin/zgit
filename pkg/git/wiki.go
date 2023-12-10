package git

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zsf-utils/idutil"
	"path/filepath"
	"strings"
	"zgit/pkg/git/command"
	"zgit/setting"
	"zgit/util"
)

const (
	WikiDefaultBranch = "master"
)

type Wiki struct {
	Id    string `json:"Id"`
	Owner User   `json:"owner"`
	Name  string `json:"name"`
	Path  string `json:"path"`
}

func InitWiki(ctx context.Context, wiki Wiki) error {
	return InitRepository(ctx, InitRepoOpts{
		Owner:         wiki.Owner,
		RepoName:      wiki.Name,
		RepoPath:      wiki.Path,
		DefaultBranch: WikiDefaultBranch,
	})
}

func UpdateWikiPage(ctx context.Context, wikiPath, pageName, content, message string) error {
	tempDir := filepath.Join(setting.TempDir(), "wiki-"+idutil.RandomUuid())
	defer util.RemoveAll(tempDir)
	hasMasterBranch, err := prepareUpdateWikiPage(ctx, wikiPath, tempDir)
	if err != nil {
		return err
	}
	object, err := HashObject(ctx, tempDir, strings.NewReader(content))
	if err != nil {
		return fmt.Errorf("hash content failed with err:%v", err)
	}
	if err = AddObjectToIndex(ctx, tempDir, DefaultFileMode, object, pageName); err != nil {
		return fmt.Errorf("addObjectToIndex failed with err:%v", err)
	}
	return afterUpdateWikiPage(ctx, tempDir, message, hasMasterBranch)
}

func DeleteWikiPage(ctx context.Context, wikiPath, pageName, message string) error {
	tempDir := filepath.Join(setting.TempDir(), "wiki-"+idutil.RandomUuid())
	defer util.RemoveAll(tempDir)
	hasMasterBranch, err := prepareUpdateWikiPage(ctx, wikiPath, tempDir)
	if err != nil {
		return err
	}
	if err = RemoveFilesFromIndex(ctx, tempDir, DefaultFileMode, pageName); err != nil {
		return fmt.Errorf("RemoveFilesFromIndex failed with err:%v", err)
	}
	return afterUpdateWikiPage(ctx, tempDir, message, hasMasterBranch)
}

func afterUpdateWikiPage(ctx context.Context, tempDir string, message string, hasMasterBranch bool) error {
	tree, err := WriteTree(ctx, tempDir)
	if err != nil {
		return fmt.Errorf("write tree failed with err:%v", err)
	}
	opts := CommitTreeOpts{
		Message: message,
	}
	if hasMasterBranch {
		opts.Parents = []string{"HEAD"}
	}
	commitHash, err := CommitTree(ctx, tempDir, tree, opts)
	if err != nil {
		return fmt.Errorf("commit tree failed with err:%v", err)
	}
	if _, err = command.NewCommand("push", DefaultRemote, fmt.Sprintf("%s:%s", commitHash, BranchPrefix+WikiDefaultBranch)).
		Run(ctx, command.WithDir(tempDir)); err != nil {
		return fmt.Errorf("push failed with err:%v", err)
	}
	return nil
}

func prepareUpdateWikiPage(ctx context.Context, wikiPath string, tempDir string) (bool, error) {
	hasMasterBranch := IsBranchExist(ctx, wikiPath, WikiDefaultBranch)
	cloneCmd := command.NewCommand("clone", "-s", "--bare", wikiPath, tempDir)
	if hasMasterBranch {
		cloneCmd.AddArgs("-b", WikiDefaultBranch)
	}
	if _, err := cloneCmd.Run(ctx); err != nil {
		return false, fmt.Errorf("clone tempDir:%s failed with err:%v", tempDir, err)
	}
	if hasMasterBranch {
		commitId, err := GetRefCommitId(ctx, tempDir, "HEAD")
		if err != nil {
			return false, fmt.Errorf("get head commitId failed with err:%v", err)
		}
		if _, err = command.NewCommand("read-tree", commitId).
			Run(ctx, command.WithDir(tempDir)); err != nil {
			return false, fmt.Errorf("read tree failed with err:%v", err)
		}
	}
	return hasMasterBranch, nil
}
