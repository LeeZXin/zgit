package cmd

import (
	"github.com/urfave/cli/v2"
	"runtime"
)

var (
	cmdList = []*cli.Command{
		Serv,
		Web,
		Proxy,
	}
)

func NewCliApp() *cli.App {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.HideHelp = true
	app.DefaultCommand = Web.Name
	app.Commands = append(app.Commands, cmdList...)
	app.Name = "zgit"
	app.Usage = "A Serv service with zsf"
	app.Description = "by default, it will start the web-server"
	app.Version = formatBuiltWith()
	return app
}

func formatBuiltWith() string {
	return " built with " + runtime.Version()
}
