package cmd

import (
	"github.com/urfave/cli/v2"
	"zgit/modules/model/clustermd"
	"zgit/pkg/sshserv/proxy"
)

var Proxy = &cli.Command{
	Name:        "proxy",
	Usage:       "This command should only be called by SSH proxy",
	Description: "Proxy provides ssh proxy for repositories",
	Action:      runProxy,
}

func runProxy(*cli.Context) error {
	// 初始化节点信息
	clustermd.InitNodesConfig()
	// 开启反向代理
	proxy.StartSSHProxy()
	return nil
}
