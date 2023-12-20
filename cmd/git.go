package cmd

import (
	"github.com/LeeZXin/zsf/starter"
	"github.com/urfave/cli/v2"
	"zgit/modules/api/hookapi"
	"zgit/modules/api/lfsapi"
	"zgit/modules/api/repoapi"
	"zgit/modules/api/sshkeyapi"
	"zgit/modules/api/userapi"
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
	// 初始化api
	lfsapi.InitApi()
	// webhook
	hookapi.InitApi()
	// 用户
	userapi.InitApi()
	// 仓库api
	repoapi.InitApi()
	// ssh公钥
	sshkeyapi.InitApi()
	starter.Run()
	return nil
}
