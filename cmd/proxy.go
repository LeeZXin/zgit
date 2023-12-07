package cmd

import (
	"github.com/urfave/cli/v2"
	"zgit/proxy"
)

var Proxy = &cli.Command{
	Name:        "proxy",
	Usage:       "This command should only be called by SSH proxy",
	Description: "Proxy provides ssh proxy for repositories",
	Action:      runProxy,
}

func runProxy(c *cli.Context) error {
	proxy.StartSSHProxy()
	return nil
}