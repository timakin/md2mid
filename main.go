package main

import (
	"os"

	"github.com/timakin/md2mid/command"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		command.InitCommand(),
		command.PublishCommand(),
	}
	app.HideVersion = true
	app.Copyright = "MIT"
	app.Usage = "Set your token to call a Medium API, and publish your markdown file."

	app.Run(os.Args)
}
