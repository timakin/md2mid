package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ericaro/frontmatter"
	medium "github.com/medium/medium-sdk-go"
	"github.com/mitchellh/go-homedir"
	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli"
)

const TokenFileName string = "~/.medium"

type Option struct {
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
	app.Usage = "Set your token to call Medium API, and publish your markdown file."

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
			err = ioutil.WriteFile(filename, []byte(token), 0644)
			if err != nil {
				return err
			}

			fmt.Printf("Token was written into %s successfully.", filename)
			return nil
		},
	}
}

func publishCommand() cli.Command {
	return cli.Command{
		Name:    "publish",
		Aliases: []string{"p"},
		Usage:   "Publish a markdown file to Medium with the status set to " + string(medium.PublishStatusDraft) + " and open the post editor page in the browser",
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
			return "", errors.New("We tried to open your browser for you automatically, but for some reason it failed. Please manually browse to https://medium.com/me/settings, generate an integration token, and use the `md2mid init <token>` command to add it.")
		}

		return "", errors.New("Could not find the token. We have opened your browser for you. Please generate an integration token in your browser window, and use `md2mid init <token>` to add it.")
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

	opts := &Option{}
	if err := frontmatter.Unmarshal(contents, opts); err != nil {
		return nil, err
	}

	postOpts := &medium.CreatePostOptions{
		Title:         opts.Title,
		ContentFormat: medium.ContentFormatMarkdown,
		Content:       opts.Content,
		Tags:          opts.Tags,
		PublishStatus: medium.PublishStatus(opts.PublishStatus),
		CanonicalURL:  opts.CanonicalURL,
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
	userId, err := getMyId()
	if err != nil {
		return err
	}

	opts.UserID = userId
	post, err := m.CreatePost(*opts)
	if err != nil {
		return err
	}

	fmt.Println("Post created. URL is ", post.URL)
	open.Run(post.URL)

	return nil
}

func getMyId() (string, error) {
	accessToken, err := getAccessToken()
	if err != nil {
		return "", err
	}

	m := medium.NewClientWithAccessToken(accessToken)
	me, err := m.GetUser("")
	if err != nil {
		return "", err
	}

	return me.ID, nil
}
