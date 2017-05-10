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

## Example

```
---
title: Sample
tags: [programming, engineer]
publish_status: draft
canonical_url: http://jamietalbot.com/posts/liverpool-fc
---

Paragraphs are separated by a blank line.

2nd paragraph. *Italic*, **bold**, and `monospace`. Itemized lists
look like:

  * this one
  * that one
  * the other one

Note that --- not considering the asterisk --- the actual text
content starts at 4-columns in.

> Block quotes are
> written like so.
>
> They can span multiple paragraphs,
> if you like.

Use 3 dashes for an em-dash. Use 2 dashes for ranges (ex., "it's all
in chapters 12--14"). Three dots ... will be converted to an ellipsis.
Unicode is supported. ☺
```

Save this file with `publish` or `p` command.
After that, console will open the page of an article.

![Sample article](https://gyazo.com/afea668308ddcd84b26d94ad0fd012d5)