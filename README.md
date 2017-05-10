md2mid
====

A CLI tool to publish your markdown files to Medium.

## Installation

`go get -u github.com/timakin/md2mid`

## Usage

```
$ md2mid -h

NAME:
   md2mid - Set your token to call a Medium API, and publish your markdown file.

USAGE:
   md2mid [global options] command [command options] [arguments...]

COMMANDS:
     init        Add an integration token to be used by this application
     publish, p  Publish a markdown file to Medium. And after the publishment, open the article page with a browser from a console.
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help

COPYRIGHT:
   MIT
```

### Init

Initialize and register an access token to access Medium API.

`md2mid init <YOUR_ACCESS_TOKEN>`

### Publish

Publish an article with your local markdown file.

`md2mid publish <MARKDOWN_FILE_PATH>`

