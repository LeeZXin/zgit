package cmd

import (
	"github.com/urfave/cli/v2"
	"runtime"
)

var (
	cmdList = []*cli.Command{
		Standalone,
		Proxy,
		Hook,
	}
)

func NewCliApp() *cli.App {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.HideHelp = true
	app.DefaultCommand = Standalone.Name
	app.Commands = append(app.Commands, cmdList...)
	app.Name = "zgit"
	app.Usage = "A Serv service with zsf"
	app.Description = "by default, it will start the git server"
	app.Version = formatBuiltWith()
	return app
}

func formatBuiltWith() string {
	return " built with " + runtime.Version()
}
