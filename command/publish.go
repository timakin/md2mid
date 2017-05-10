package command

import (
	"log"

	medium "github.com/medium/medium-sdk-go"
	"github.com/skratchdot/open-golang/open"
	"github.com/timakin/md2mid/util"
	"github.com/urfave/cli"
)

func PublishCommand() cli.Command {
	return cli.Command{
		Name:    "publish",
		Aliases: []string{"p"},
		Usage:   "Publish a markdown file to Medium. And after the publishment, open the article page with a browser from a console.",
		Action:  publish,
	}
}

func publish(c *cli.Context) error {
	opts, err := util.ParseOpts(c.Args().First())
	if err != nil {
		return err
	}

	accessToken, err := util.GetAccessToken()
	if err != nil {
		return err
	}

	m := medium.NewClientWithAccessToken(accessToken)

	userId, err := util.GetMyId(m)
	if err != nil {
		return err
	}

	opts.UserID = userId
	post, err := m.CreatePost(*opts)
	if err != nil {
		return err
	}

	log.Printf("Successfully generated your post! Check it out on %s", post.URL)
	open.Run(post.URL)

	return nil
}
