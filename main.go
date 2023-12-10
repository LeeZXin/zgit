package main

import (
	"fmt"
	"os"
	"zgit/cmd"
)

func main() {
	app := cmd.NewCliApp()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
