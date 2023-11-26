package main

import (
	"github.com/LeeZXin/zsf/logger"
	"os"
	"zgit/cmd"
)

func main() {
	app := cmd.NewCliApp()
	err := app.Run(os.Args)
	if err != nil {
		logger.Logger.Panic(err)
	}
}
