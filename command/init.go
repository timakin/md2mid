package command

import (
	"errors"

	"github.com/timakin/md2mid/util"
	"github.com/urfave/cli"
)

func InitCommand() cli.Command {
	return cli.Command{
		Name:    "init",
		Aliases: []string{},
		Usage:   "Add an integration token to be used by this application",
		Action: func(c *cli.Context) error {
			token := c.Args().First()
			if token == "" {
				return errors.New("Input your token as an arguement.")
			}
			err := util.WriteAccessToken(token)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
