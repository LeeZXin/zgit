package cmd

import (
	"github.com/LeeZXin/zsf/starter"
	"github.com/urfave/cli/v2"
	"zgit/api/lfsapi"
	"zgit/pkg/git"
	"zgit/pkg/sshserv"
)

var Git = &cli.Command{
	Name:        "git",
	Usage:       "This command starts zgit git server",
	Description: "zgit",
	Action:      runGit,
}

func runGit(*cli.Context) error {
	sshserv.InitSsh()
	git.InitGit()
	lfsapi.InitLfsHttpApi()
	starter.Run()
	return nil
}
