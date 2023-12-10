package git

import (
	"context"
	"zgit/pkg/git/command"
)

func UpdateServerInfo(ctx context.Context, repoPath string) error {
	_, err := command.NewCommand("update-server-info").Run(ctx, command.WithDir(repoPath))
	return err
}

func RevParse(ctx context.Context, repoPath string, args ...string) error {
	_, err := command.NewCommand("rev-parse").AddArgs(args...).Run(ctx, command.WithDir(repoPath))
	return err
}
