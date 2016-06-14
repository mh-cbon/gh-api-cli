# gh-api-cli

Command line client for github api

# Install

```sh
mkdir -p $GOPATH/github.com/mh-cbon
cd $GOPATH/github.com/mh-cbon
git clone https://github.com/mh-cbon/gh-api-cli.git
cd gh-api-cli
glide install
go install
```

# Usage

```sh
NAME:
   gh-api-cli - Github api command line client

USAGE:
   gh-api-cli <cmd> <options>

VERSION:
   0.0.0

COMMANDS:
     add-auth     Add a new authorization
     list-auth	  List authorizations
     rm-auth	    Remove an authorization
     get-auth	    Get token from a locally saved authorization

GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
```

#### add-auth
```sh
NAME:
   gh-api-cli add-auth - Add a new authorization

USAGE:
   gh-api-cli add-auth [command options] [arguments...]

OPTIONS:
   --username value, -u value	Github username
   --password value, -p value	Github password
   --name value, -n value	Name of the authorization
   --rights value, -r value	Permissions to set
```

#### list-auth
```sh
NAME:
   gh-api-cli list-auth - List authorizations

USAGE:
   gh-api-cli list-auth [command options] [arguments...]

OPTIONS:
   --username value, -u value	Github username
   --password value, -p value	Github password
   --name value, -n value	Name of the authorization
```

#### rm-auth
```sh
NAME:
   gh-api-cli rm-auth - Remove an existing authorization

USAGE:
   gh-api-cli rm-auth [command options] [arguments...]

OPTIONS:
   --username value, -u value	Github username
   --password value, -p value	Github password
   --name value, -n value	Name of the authorization
```

#### get-auth
```sh
NAME:
   gh-api-cli get-auth - Get token from a locally saved authorization

USAGE:
   gh-api-cli get-auth [command options] [arguments...]

OPTIONS:
   --name value, -n value	Name of the authorization
```
