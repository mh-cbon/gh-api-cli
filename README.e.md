---
License: MIT
LicenseFile: LICENSE
LicenseColor: yellow
---
# {{.Name}}

{{template "badge/travis" .}} {{template "badge/appveyor" .}} {{template "badge/goreport" .}} {{template "badge/godoc" .}} {{template "license/shields" .}}

{{pkgdoc}}

This tool is part of the [go-github-release workflow](https://github.com/mh-cbon/go-github-release)

# {{toc 5}}

# Install
{{template "gh/releases" .}}

#### Glide
{{template "glide/install" .}}

#### Chocolatey
{{template "choco/install" .}}

#### linux rpm/deb repository
{{template "linux/gh_src_repo" .}}

#### linux rpm/deb standalone package
{{template "linux/gh_pkg" .}}

# Cli

###### {{exec "gh-api-cli" "-help" | color "sh"}}

###### {{exec "gh-api-cli" "add-auth" "-help" | color "sh"}}

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


###### {{exec "gh-api-cli" "list-auth" "-help" | color "sh"}}

```sh
EXAMPLE
  gh-api-cli list-auth -n test # will prompt for username/password
  gh-api-cli list-auth -n test -u your -p pwd # won t prompt unless you have 2F ident on
```

###### {{exec "gh-api-cli" "rm-auth" "-help" | color "sh"}}

```
EXAMPLE
  gh-api-cli rm-auth -n test # will prompt for username/password
  gh-api-cli rm-auth -n test -u your -p pwd # won t prompt unless you have 2F ident on
```

###### {{exec "gh-api-cli" "get-auth" "-help" | color "sh"}}

```
EXAMPLE
  gh-api-cli get-auth -n test
```

###### {{exec "gh-api-cli" "create-release" "-help" | color "sh"}}

```
EXAMPLE
  gh-api-cli create-release -n test -o mh-cbon -r gh-api-cli --ver 0.0.1
  gh-api-cli create-release -n test --guess --ver 0.0.1
```

###### {{exec "gh-api-cli" "rm-release" "-help" | color "sh"}}

```
EXAMPLE
  gh-api-cli rm-release -n test -o mh-cbon -r gh-api-cli --ver 0.0.1
  gh-api-cli rm-release -n test --guess --ver 0.0.1
```

###### {{exec "gh-api-cli" "upload-release-asset" "-help" | color "sh"}}

```
EXAMPLE
  gh-api-cli upload-release-asset -n test -g README.md -o mh-cbon -r gh-api-cli --ver 0.0.1
  gh-api-cli upload-release-asset -n test -g README.md --guess --ver 0.0.1
```

###### {{exec "gh-api-cli" "rm-assets" "-help" | color "sh"}}

```
EXAMPLE
  gh-api-cli rm-assets -n test --glob file.package -o mh-cbon -r gh-api-cli --ver 0.0.1
  gh-api-cli rm-assets -n test -g file.package --guess --ver 0.0.1
```

###### {{exec "gh-api-cli" "dl-assets" "-help" | color "sh"}}

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
