package cmd

import (
	"github.com/urfave/cli/v2"
	"zgit/git"
	"zgit/ssh"
)

var Git = &cli.Command{
	Name:        "git",
	Usage:       "This command starts zgit git server",
	Description: "zgit",
	Action:      runGit,
}

func runGit(c *cli.Context) error {
	ssh.InitSsh()
	git.InitGit()
	ssh.StartServer()
	return nil
}
