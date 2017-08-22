# gh-api-cli

[![travis Status](https://travis-ci.org/mh-cbon/gh-api-cli.svg?branch=master)](https://travis-ci.org/mh-cbon/gh-api-cli) [![Appveyor Status](https://ci.appveyor.com/api/projects/status/github/mh-cbon/gh-api-cli?branch=master&svg=true)](https://ci.appveyor.com/projects/mh-cbon/gh-api-cli) [![Go Report Card](https://goreportcard.com/badge/github.com/mh-cbon/gh-api-cli)](https://goreportcard.com/report/github.com/mh-cbon/gh-api-cli) [![GoDoc](https://godoc.org/github.com/mh-cbon/gh-api-cli?status.svg)](http://godoc.org/github.com/mh-cbon/gh-api-cli) [![MIT License](http://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Package gh-api-cli is a command line utility to work with github api.


This tool is part of the [go-github-release workflow](https://github.com/mh-cbon/go-github-release)

# TOC
- [Install](#install)
  - [Glide](#glide)
  - [Chocolatey](#chocolatey)
  - [linux rpm/deb repository](#linux-rpmdeb-repository)
  - [linux rpm/deb standalone package](#linux-rpmdeb-standalone-package)
- [Cli](#cli)
- [Notes](#notes)
- [Todo](#todo)
- [Recipes](#recipes)
  - [Testing](#testing)
  - [Release the project](#release-the-project)
- [History](#history)

# Install
Check the [release page](https://github.com/mh-cbon/gh-api-cli/releases)!

#### Glide
```sh
mkdir -p $GOPATH/src/github.com/mh-cbon/gh-api-cli
cd $GOPATH/src/github.com/mh-cbon/gh-api-cli
git clone https://github.com/mh-cbon/gh-api-cli.git .
glide install
go install
```

#### Chocolatey
```sh
choco install gh-api-cli
```

#### linux rpm/deb repository
```sh
wget -O - https://raw.githubusercontent.com/mh-cbon/latest/master/source.sh \
| GH=mh-cbon/gh-api-cli sh -xe
# or
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/source.sh \
| GH=mh-cbon/gh-api-cli sh -xe
```

#### linux rpm/deb standalone package
```sh
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/gh-api-cli sh -xe
# or
wget -q -O - --no-check-certificate \
https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/gh-api-cli sh -xe
```

# Cli

###### gh-api-cli -help
```sh
NAME:
   gh-api-cli - Github api command line client

USAGE:
   gh-api-cli <cmd> <options>

VERSION:
   0.0.0

COMMANDS:
     add-auth              Add a new authorization
     list-auth             List authorizations
     rm-auth               Remove an authorization
     get-auth              Get token from a locally saved authorization
     create-release        Create a release
     rm-release            Delete a release
     upload-release-asset  Upload assets to a release
     dl-assets             Download assets
     rm-assets             Delete assets
     help, h               Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

###### gh-api-cli add-auth -help
```sh
NAME:
   gh-api-cli add-auth - Add a new authorization

USAGE:
   gh-api-cli add-auth [command options] [arguments...]

OPTIONS:
   --username value, -u value  Github username
   --password value, -p value  Github password
   --name value, -n value      Name of the authorization to create
   --rights value, -r value    Permissions to set
```

```sh
EXAMPLE
  gh-api-cli add-auth -n test -r user -r repo # will prompt for username/password
  gh-api-cli add-auth -n test -r user -u your -p pwd # won t prompt unless you have 2F ident on
  gh-api-cli add-auth -n test -r user,repo -u your -p pwd
```

Where `rights` contains some of :

| Name | - | - |
| --- | --- | --- |
| user | user:email | user:follow |
| public_repo | repo | repo_deployment |
| notifications | repo:status | delete_repo |
| gist |read:repo_hook | write:repo_hook |
| admin:org_hook | read:org | write:org |
| admin | admin:org | admin:repo_hook |
| admin:public_key | read:public_key | write:public_key |
| read:gpg_key | write:gpg_key | admin:gpg_key |


###### gh-api-cli list-auth -help
```sh
NAME:
   gh-api-cli list-auth - List authorizations

USAGE:
   gh-api-cli list-auth [command options] [arguments...]

OPTIONS:
   --username value, -u value  Github username
   --password value, -p value  Github password
```

```sh
EXAMPLE
  gh-api-cli list-auth -n test # will prompt for username/password
  gh-api-cli list-auth -n test -u your -p pwd # won t prompt unless you have 2F ident on
```

###### gh-api-cli rm-auth -help
```sh
NAME:
   gh-api-cli rm-auth - Remove an authorization

USAGE:
   gh-api-cli rm-auth [command options] [arguments...]

OPTIONS:
   --username value, -u value  Github username
   --password value, -p value  Github password
   --name value, -n value      Name of the authorization to delete
```

```
EXAMPLE
  gh-api-cli rm-auth -n test # will prompt for username/password
  gh-api-cli rm-auth -n test -u your -p pwd # won t prompt unless you have 2F ident on
```

###### gh-api-cli get-auth -help
```sh
NAME:
   gh-api-cli get-auth - Get token from a locally saved authorization

USAGE:
   gh-api-cli get-auth [command options] [arguments...]

OPTIONS:
   --name value, -n value  Name of the authorization
```

```
EXAMPLE
  gh-api-cli get-auth -n test
```

###### gh-api-cli create-release -help
```sh
NAME:
   gh-api-cli create-release - Create a release

USAGE:
   gh-api-cli create-release [command options] [arguments...]

OPTIONS:
   --name value, -n value        Name of the authorization to use for identification
   --token value, -t value       Value of a personal access token
   --owner value, -o value       Repository owner
   --repository value, -r value  Repository name
   --guess                       Guess repository and user name from the cwd
   --ver value                   Version name
   --author value, -a value      Release author github username
   --draft value, -d value       Make a draft release (default: "no")
   --changelog value, -c value   A command to generate the description body of the release
```

```
EXAMPLE
  gh-api-cli create-release -n test -o mh-cbon -r gh-api-cli --ver 0.0.1
  gh-api-cli create-release -n test --guess --ver 0.0.1
```

###### gh-api-cli rm-release -help
```sh
NAME:
   gh-api-cli rm-release - Delete a release

USAGE:
   gh-api-cli rm-release [command options] [arguments...]

OPTIONS:
   --name value, -n value        Name of the authorization to use for identification
   --token value, -t value       Value of a personal access token
   --owner value, -o value       Repository owner
   --repository value, -r value  Repository name
   --guess                       Guess repository and user name from the cwd
   --ver value                   Version name
```

```
EXAMPLE
  gh-api-cli create-release -n test -o mh-cbon -r gh-api-cli --ver 0.0.1
  gh-api-cli create-release -n test --guess --ver 0.0.1
```

###### gh-api-cli upload-release-asset -help
```sh
NAME:
   gh-api-cli upload-release-asset - Upload assets to a release

USAGE:
   gh-api-cli upload-release-asset [command options] [arguments...]

OPTIONS:
   --name value, -n value        Name of the authorization to use for identification
   --token value, -t value       Value of a personal access token
   --glob value, -g value        Glob pattern of files to upload
   --owner value, -o value       Repository owner
   --repository value, -r value  Repository name
   --guess                       Guess repository and user name from the cwd
   --ver value                   Version name
```

```
EXAMPLE
  gh-api-cli upload-release-asset -n test -g README.md -o mh-cbon -r gh-api-cli --ver 0.0.1
  gh-api-cli upload-release-asset -n test -g README.md --guess --ver 0.0.1
```

###### gh-api-cli rm-assets -help
```sh
NAME:
   gh-api-cli rm-assets - Delete assets

USAGE:
   gh-api-cli rm-assets [command options] [arguments...]

OPTIONS:
   --name value, -n value        Name of the authorization to use for identification
   --token value, -t value       Value of a personal access token
   --glob value, -g value        Glob pattern of files to download
   --owner value, -o value       Repository owner
   --repository value, -r value  Repository name
   --guess                       Guess repository and user name from the cwd
   --ver value                   Version constraint
```

```
EXAMPLE
  gh-api-cli upload-release-asset -n test -g README.md -o mh-cbon -r gh-api-cli --ver 0.0.1
  gh-api-cli upload-release-asset -n test -g README.md --guess --ver 0.0.1
```

###### gh-api-cli dl-assets -help
```sh
NAME:
   gh-api-cli dl-assets - Download assets

USAGE:
   gh-api-cli dl-assets [command options] [arguments...]

OPTIONS:
   --name value, -n value        Name of the authorization to use for identification
   --token value, -t value       Value of a personal access token
   --glob value, -g value        Glob pattern of files to download
   --out value                   Out format to write files (default: "%f")
   --owner value, -o value       Repository owner
   --repository value, -r value  Repository name
   --guess                       Guess repository and user name from the cwd
   --ver value                   Version constraint
   --skip-prerelease value       Skip prerelease releases (yes|no) (default: "no")
```

```
EXAMPLE
  gh-api-cli dl-assets -o mh-cbon -r gh-api-cli --ver 4.x --out "dl/%r.%v-%a.%e"
  gh-api-cli dl-assets -o mh-cbon -r gh-api-cli --ver 0.0.1 --out dl/%f
  gh-api-cli dl-assets -o mh-cbon -r gh-api-cli --ver 0.0.1 --out dl/%f -g '*amd64*deb'
  gh-api-cli dl-assets -o mh-cbon -r gh-api-cli --ver latest --out dl/%s/%r.%v-%a.%e
  gh-api-cli dl-assets -o mh-cbon -r gh-api-cli --out "dl/%s/%r-%v-%a.%e" --ver ">0.0.10"
```

# Notes

When you `add, remove, list` authorizations, personal access token authentication is not permitted, [see this](https://developer.github.com/v3/oauth_authorizations/#deprecation-notice)

You are required to use a password.

# Todo

- add a command to clean up old gh releases,
something that would help to keep only N most recent releases for each major version.

# Recipes

#### Testing

```sh
 (USER=xxx PWRD=yyy ./test.sh | grep "OK, ALL FINE") || (echo "" && echo "" && echo "beep boop failed")
```

#### Release the project

```sh
gump patch -d # check
gump patch # bump
```

# History

[CHANGELOG](CHANGELOG.md)
