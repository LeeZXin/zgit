package reposrv

import (
	"context"
	"zgit/modules/model/repomd"
)

func GetRepoInfoByRelativePath(ctx context.Context, path string) (repomd.RepoInfo, bool, error) {
	return repomd.RepoInfo{}, true, nil
}
