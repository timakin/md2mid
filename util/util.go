package util

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/ericaro/frontmatter"
	medium "github.com/medium/medium-sdk-go"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/skratchdot/open-golang/open"
)

type FrontmatterOption struct {
	Title         string   `fm:"title"`
	Tags          []string `fm:"tags"`
	Content       string   `fm:"content"`
	PublishStatus string   `fm:"publishstatus"`
	CanonicalURL  string   `fm:"canonicalurl"`
}

const TokenFileName string = "~/.medium"

func WriteAccessToken(token string) error {
	filename, err := homedir.Expand(TokenFileName)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, []byte(token), 0644)
	if err != nil {
		return err
	}

	log.Printf("Successfully the token was saved into %s.", filename)
	return nil
}

func GetAccessToken() (string, error) {
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

func ParseOpts(filename string) (*medium.CreatePostOptions, error) {
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

func GetMyId(m *medium.Medium) (string, error) {
	me, err := m.GetUser("")
	if err != nil {
		return "", err
	}

	return me.ID, nil
}
