package git

import (
	"context"
	"zgit/git/command"
)

const (
	DefaultRemote = "origin"
)

func AddRemote(ctx context.Context, repoPath, name, url string, tryFetch bool) error {
	cmd := command.NewCommand("remote", "add", name, url)
	if tryFetch {
		cmd.AddArgs("-f")
	}
	_, err := cmd.Run(ctx, command.WithDir(repoPath))
	return err
}

func RemoveRemote(ctx context.Context, repoPath, name string) error {
	_, err := command.NewCommand("remote", "rm", name).Run(ctx, command.WithDir(repoPath))
	return err
}
