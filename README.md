# gh-api-cli

Command line client for github api

# Install

You can grab a pre-built binary file in the [releases page](https://github.com/mh-cbon/gh-api-cli/releases)

```sh
mkdir -p $GOPATH/github.com/mh-cbon
cd $GOPATH/github.com/mh-cbon
git clone https://github.com/mh-cbon/gh-api-cli.git
cd gh-api-cli
glide install
go install
```

# Usage

```
NAME:
   gh-api-cli - Github api command line client

USAGE:
   gh-api-cli <cmd> <options>

VERSION:
   0.0.0

COMMANDS:
     add-auth                 Add a new authorization
     list-auth                List authorizations
     rm-auth                  Remove an authorization
     get-auth                 Get token from a locally saved authorization
     create-release           Create a new release
     upload-release-asset     Upload assets to a release

GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
```

#### add-auth
```
NAME:
   gh-api-cli add-auth - Add a new authorization

USAGE:
   gh-api-cli add-auth [command options] [arguments...]

OPTIONS:
   --username value, -u value    Github username
   --password value, -p value    Github password
   --name value, -n value        Name of the authorization
   --rights value, -r value      Permissions to set
```

```
EXAMPLE
  gh-api-cli add-auth -n test -r user -r repo # will prompt for username/password
  gh-api-cli add-auth -n test -r user -u your -p pwd # won t prompt unless you have 2F ident on
```

Where `rights` is one of :

```
user              user:email
user:follow       public_repo
repo              repo_deployment
repo:status       delete_repo
notifications     gist
read:repo_hook    write:repo_hook
admin:repo_hook   admin:org_hook
admin             read:org
write:org         admin:org
read:public_key   write:public_key
admin:public_key  read:gpg_key
write:gpg_key     admin:gpg_key
```

#### list-auth
```
NAME:
   gh-api-cli list-auth - List authorizations

USAGE:
   gh-api-cli list-auth [command options] [arguments...]

OPTIONS:
   --username value, -u value   Github username
   --password value, -p value   Github password
   --name value, -n value       Name of the authorization
```

```
EXAMPLE
  gh-api-cli list-auth -n test # will prompt for username/password
  gh-api-cli list-auth -n test -u your -p pwd # won t prompt unless you have 2F ident on
```

#### rm-auth
```
NAME:
   gh-api-cli rm-auth - Remove an existing authorization

USAGE:
   gh-api-cli rm-auth [command options] [arguments...]

OPTIONS:
   --username value, -u value   Github username
   --password value, -p value   Github password
   --name value, -n value       Name of the authorization
```

```
EXAMPLE
  gh-api-cli rm-auth -n test # will prompt for username/password
  gh-api-cli rm-auth -n test -u your -p pwd # won t prompt unless you have 2F ident on
```

#### get-auth
```
NAME:
   gh-api-cli get-auth - Get token from a locally saved authorization

USAGE:
   gh-api-cli get-auth [command options] [arguments...]

OPTIONS:
   --name value, -n value       Name of the authorization
```

```
EXAMPLE
  gh-api-cli get-auth -n test
```

#### upload-release-asset
```
NAME:
   gh-api-cli upload-release-asset - Upload assets to a release

USAGE:
   gh-api-cli upload-release-asset [command options] [arguments...]

OPTIONS:
   --name value, -n value          Name of the authorization
   --glob value, -g value          Glob pattern of files to upload
   --owner value, -o value         Repo owner
   --repository value, -r value    Repo name
   --ver value                     Version name
```

```
EXAMPLE
  gh-api-cli upload-release-asset -n test -g README.md -o mh-cbon -r gh-api-cli --ver 0.0.1
```

#### create-release
```
NAME:
   gh-api-cli create-release - Create a release

USAGE:
   gh-api-cli create-release [command options] [arguments...]

OPTIONS:
   --name value, -n value          Name of the authorization
   --owner value, -o value         Repo owner
   --repository value, -r value    Repo name
   --ver value                     Version name
   --author value, -a value        Release author name
   --email value, -e value         Release author email
   --draft, -d                     Make a draft release

```

```
EXAMPLE
  gh-api-cli create-release -n test -o mh-cbon -r gh-api-cli --ver 0.0.1
```
