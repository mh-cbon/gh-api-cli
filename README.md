# gh-api-cli

Command line client for github api

This tool is part of the [go-github-release workflow](https://github.com/mh-cbon/go-github-release)

# Install

Pick an msi package [here](https://github.com/mh-cbon/gh-api-cli/releases)!

__chocolatey__

```sh
choco install gh-api-cli
```

__deb/ubuntu/rpm repositories__

```sh
wget -O - https://raw.githubusercontent.com/mh-cbon/latest/master/source.sh \
| GH=mh-cbon/gh-api-cli sh -xe
# or
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/source.sh \
| GH=mh-cbon/gh-api-cli sh -xe
```

__deb/ubuntu/rpm package__

```sh
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/gh-api-cli sh -xe
# or
wget -q -O - --no-check-certificate \
https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/gh-api-cli sh -xe
```

__go__

```sh
mkdir -p $GOPATH/src/github.com/mh-cbon
cd $GOPATH/src/github.com/mh-cbon
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
     rm-release               Delete a release
     upload-release-asset     Upload assets to a release
     dl-assets                Download assets
     rm-assets                Delete assets

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
   --name value, -n value        Name of the authorization to add
   --rights value, -r value      Permissions to set
```

```
EXAMPLE
  gh-api-cli add-auth -n test -r user -r repo # will prompt for username/password
  gh-api-cli add-auth -n test -r user -u your -p pwd # won t prompt unless you have 2F ident on
  gh-api-cli add-auth -n test -r user,repo -u your -p pwd
```

Where `rights` contains some of :

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
   --name value, -n value       Name of the authorization to delete
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
   --name value, -n value       Name of the authorization to look for
```

```
EXAMPLE
  gh-api-cli get-auth -n test
```

#### create-release
```
NAME:
   gh-api-cli create-release - Create a release

USAGE:
   gh-api-cli create-release [command options] [arguments...]

OPTIONS:
   --name value, n value           Name of the locally saved authentication.
   --token value, t value          Personal access token authentication.
   --owner value, -o value         Repo owner
   --repository value, -r value    Repo name
   --ver value                     Version name
   --author value, -a value        Release author name
   --draft value, -d value         Make a draft release, value=yes|1|true|no|0|false
   --changelog cmd, -c cmd         A command to generate body content of the release
```

```
EXAMPLE
  gh-api-cli create-release -n test -o mh-cbon -r gh-api-cli --ver 0.0.1
```

#### rm-release
```
NAME:
   gh-api-cli rm-release - Delete a release

USAGE:
   gh-api-cli rm-release [command options] [arguments...]

OPTIONS:
   --name value, n value           Name of the locally saved authentication.
   --token value, t value          Personal access token authentication.
   --owner value, -o value         Repo owner
   --repository value, -r value    Repo name
   --ver value                     Version name
```

```
EXAMPLE
  gh-api-cli create-release -n test -o mh-cbon -r gh-api-cli --ver 0.0.1
```

#### upload-release-asset
```
NAME:
   gh-api-cli upload-release-asset - Upload assets to a release

USAGE:
   gh-api-cli upload-release-asset [command options] [arguments...]

OPTIONS:
   --name value, n value           Name of the locally saved authentication.
   --token value, t value          Personal access token authentication.
   --glob value, -g value          Glob pattern of files to upload
   --owner value, -o value         Repo owner
   --repository value, -r value    Repo name
   --ver value                     Version name
```

```
EXAMPLE
  gh-api-cli upload-release-asset -n test -g README.md -o mh-cbon -r gh-api-cli --ver 0.0.1
```

#### rm-assets
```
NAME:
   gh-api-cli rm-assets - Delete assets of a release

USAGE:
   gh-api-cli rm-assets [command options] [arguments...]

OPTIONS:
   --name value, n value           Name of the locally saved authentication.
   --token value, t value          Personal access token authentication.
   --glob value, -g value          Glob pattern of files to upload
   --owner value, -o value         Repo owner
   --repository value, -r value    Repo name
   --ver value                     Version name
```

```
EXAMPLE
  gh-api-cli upload-release-asset -n test -g README.md -o mh-cbon -r gh-api-cli --ver 0.0.1
```

#### dl-assets
```
NAME:
   gh-api-cli dl-assets - Download assets

USAGE:
   gh-api-cli dl-assets [command options] [arguments...]

OPTIONS:
   --name value, n value           Name of the locally saved authentication.
   --token value, t value          Personal access token authentication.
   --owner value, -o value         Repo owner
   --repository value, -r value    Repo name
   --glob value, -g value          A glob to match files to download.
                                   It resolves to a regexp like '(i?)^glob$'.
                                   Stars '*' are replace by '.+'.
   --skip-prerelease yes|no        if yes, skips pre-releases from the selection.
   --version constraint            A version constraint,
                                   Special value 'latest' is acceptable.
   --out value                     A formatted string to write files.
                                   It can contain token such as
                                   %f: full filename
                                   %o: repository owner
                                   %r: repository name
                                   %e: file extension, minus dot prefix, detected JIT
                                   %s: target system (windows, darwin, linux), detected JIT
                                   %a: architecture (amd64, 386), detected JIT
                                   %v: version the asset is attached to
```

```
EXAMPLE
  gh-api-cli dl-assets -o mh-cbon -r gh-api-cli --ver 0.0.1
  gh-api-cli dl-assets -o mh-cbon -r gh-api-cli --ver 0.0.1 --out dl/%f
  gh-api-cli dl-assets -o mh-cbon -r gh-api-cli --ver 0.0.1 --out dl/%f -g '*amd64*deb'
  gh-api-cli dl-assets -o mh-cbon -r gh-api-cli --ver latest --out dl/%s/%r.%v-%a.%e
  gh-api-cli dl-assets -o mh-cbon -r gh-api-cli --out "dl/%s/%r-%v-%a.%e" --ver ">0.0.10"
```

# Notes

When you `add, remove, list` authorizations, personal access token authentication is not permitted, [see this](https://developer.github.com/v3/oauth_authorizations/#deprecation-notice)

# Testing

```sh
 (USER=xxx PWRD=yyy ./test.sh | grep "OK, ALL FINE") || (echo "" && echo "" && echo "beep boop failed")
```

# Todo

- add a command to clean up old gh releases,
something that would help to keep only N most recent releases for each major version.
