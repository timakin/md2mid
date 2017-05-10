package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/ericaro/frontmatter"
	medium "github.com/medium/medium-sdk-go"
	"github.com/mitchellh/go-homedir"
	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli"
)

const TokenFileName string = "~/.medium"

type FrontmatterOption struct {
	Title         string   `fm:"title"`
	Tags          []string `fm:"tags"`
	Content       string   `fm:"content"`
	PublishStatus string   `fm:"publish_status"`
	CanonicalURL  string   `fm:"canonical_url"`
}

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		initCommand(),
		publishCommand(),
	}
	app.HideVersion = true
	app.Copyright = "MIT"
	app.Usage = "Set your token to call a Medium API, and publish your markdown file."

	app.Run(os.Args)
}

func initCommand() cli.Command {
	return cli.Command{
		Name:    "init",
		Aliases: []string{},
		Usage:   "Add an integration token to be used by this application",
		Action: func(c *cli.Context) error {
			filename, err := homedir.Expand(TokenFileName)
			if err != nil {
				return err
			}

			token := c.Args().First()
			if token == "" {
				return errors.New("Input your token as an arguement.")
			}
			err = ioutil.WriteFile(filename, []byte(token), 0644)
			if err != nil {
				return err
			}

			log.Printf("Successfully the token was saved into %s.", filename)
			return nil
		},
	}
}

func publishCommand() cli.Command {
	return cli.Command{
		Name:    "publish",
		Aliases: []string{"p"},
		Usage:   "Publish a markdown file to Medium. And after the publishment, open the article page with a browser from a console.",
		Action:  publish,
	}
}

func getAccessToken() (string, error) {
	filename, err := homedir.Expand(TokenFileName)
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadFile(filename)

	// If the file is not found, it means we have no tokens. Try and make it more user-friendly.
	if err != nil && os.IsNotExist(err) {
		err := open.Run("https://medium.com/me/settings")
		if err != nil {
			return "", errors.New("Failed to open `https://medium.com/me/settings`")
		}

		return "", errors.New("Failed to parse your token. Please `md2mid init` at first to initialize this tool's setting.")
	}

	// It's not a file not found error, so just return it
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func parseOpts(filename string) (*medium.CreatePostOptions, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	fmOpts := &FrontmatterOption{}
	if err := frontmatter.Unmarshal(contents, fmOpts); err != nil {
		return nil, err
	}

	postOpts := &medium.CreatePostOptions{
		Title:         fmOpts.Title,
		ContentFormat: medium.ContentFormatMarkdown,
		Content:       fmOpts.Content,
		Tags:          fmOpts.Tags,
		PublishStatus: medium.PublishStatus(fmOpts.PublishStatus),
		CanonicalURL:  fmOpts.CanonicalURL,
	}

	return postOpts, nil
}

func publish(c *cli.Context) error {
	opts, err := parseOpts(c.Args().First())
	if err != nil {
		return err
	}

	accessToken, err := getAccessToken()
	if err != nil {
		return err
	}

	m := medium.NewClientWithAccessToken(accessToken)

	userId, err := getMyId(m)
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

func getMyId(m *medium.Medium) (string, error) {
	me, err := m.GetUser("")
	if err != nil {
		return "", err
	}

	return me.ID, nil
}
