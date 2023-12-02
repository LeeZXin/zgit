package cmd

import (
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
	starter.Run(options...)
	return nil
}
