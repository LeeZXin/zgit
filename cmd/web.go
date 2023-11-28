package cmd

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zsf/starter"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/urfave/cli/v2"
	"zgit/git"
	"zgit/ssh"
)

var Web = &cli.Command{
	Name:        "web",
	Usage:       "This command starts zgit server",
	Description: "zgit",
	Action:      runWeb,
}

func runWeb(c *cli.Context) error {
	options := make([]zsf.Option, 0)
	args := c.Args().Slice()
	for _, arg := range args {
		if arg == "debug" {
			options = append(options, zsf.WithRunMode("debug"))
		}
	}
	ssh.InitSsh()
	git.InitGit()
	_, err := git.GetGitLogCommitList(context.Background(), "/Users/lizexin/go/src/zgit/data/repo/bb/runner-test.git", "9727e46b0db7a5448d6e2aa7b39210ab669d0674", "refs/heads/uat", true)
	fmt.Println(err)
	starter.Run(options...)
	return nil
}
