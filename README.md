# Git Browse

> A handy command line utility for generating web links to files and dirs in the upstream repo.

## Features

Natively Supports:

  * Github
  * Sourcegraph

For all other web accessible repository browsers it can be extended using configuration.

## Usage

First you need to tell this tool where your repo is located. You can refer to the configuration:

```
$ git-browse config --format=cli

git config browse.default "github"

git config browse.service.github.base "https://github.com"
git config browse.service.github.repo ""
git config browse.service.github.url-file "{{.Base}}/{{.Repo}}/blob/{{.Ref}}/{{.Path}}{{if .WithLine}}#L{{.LineS}}{{else if .WithRange}}#L{{.LineS}}-L{{.LineE}}{{end}}"
git config browse.service.github.url-file "{{.Base}}/{{.Repo}}/tree/{{.Ref}}/{{.Path}}"

git config browse.service.sourcegraph.base "https://sourcegraph.com"
git config browse.service.sourcegraph.repo ""
git config browse.service.sourcegraph.url-file "{{.Base}}/{{.Repo}}/{{.Ref}}/blob/{{.Path}}{{if .WithLine}}#L{{.LineS}}{{else if .WithRange}}#L{{.LineS}}-{{.LineE}}{{end}}"
git config browse.service.sourcegraph.url-file "{{.Base}}/{{.Repo}}/{{.Ref}}/tree/{{.Path}}"
```

If your repository is on Github, you need to set `browse.service.github.repo` git config option to the fully qualified name of your repo.

Using this repo as an example:

```
git config browse.service.github.repo NonLogicalDev/util.git-browse
```

Once you have configured this tool, you can easily fetch a link to any dir or file by using it like so:

```
$ git browse url ./cmd/git-browse/main.go

https://github.com/NonLogicalDev/util.git-browse/blob/HEAD/cmd/git-browse/main.go
```

