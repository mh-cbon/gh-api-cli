# Changelog - gh-api-cli

### 4.0.3

__Changes__

- ci: upate tokens
- asset download: fix concurrency

__Contributors__

- mh-cbon

Released by mh-cbon, Wed 23 Aug 2017 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/4.0.2...4.0.3#diff)
______________

### 4.0.2

__Changes__

- bump: fix README generation
- asset download: add concurrency, improve output

__Contributors__

- mh-cbon

Released by mh-cbon, Fri 21 Apr 2017 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/4.0.1...4.0.2#diff)
______________

### 4.0.1

__Changes__

- appveyor: tryfix to build msi package

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 05 Jan 2017 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/4.0.0...4.0.1#diff)
______________

### 4.0.0

__Changes__

- cli(minor): add new guess parameter to compute repository owner/name from the cwd
- cli(break): remove author email parameter from create-release command
- gh.GetReleaseById(break): renamed to GetReleaseByID
- gh.ReleaseId(break): renamed to ReleaseID
- gh.CreateRelease(break): removed authoremail parameter
- dl.Asset(break): renamed SourceUrl to SourceURL
- Applied linters
- Updated dependencies
- appveyor: update choco key

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 05 Jan 2017 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/3.0.4...4.0.0#diff)
______________

### 3.0.4

__Changes__

- choco: refine the choco package information to pass chocolatey validation and get th package published
- README
- appveyor: update gh token
- release: update release script
- README: update install section
- build: update travis file to use more env
- build: changed appveyor gh token
- changelog: 3.0.3

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 04 Aug 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/3.0.3...3.0.4#diff)
______________

### 3.0.3

__Changes__

- build: fix appveyor script

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 30 Jul 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/3.0.2...3.0.3#diff)
______________

### 3.0.2

__Changes__

- build: fix appveyor script

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 30 Jul 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/3.0.1...3.0.2#diff)
______________

### 3.0.1

__Changes__

- build: add choco package
- README: updated usage section

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 30 Jul 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/3.0.0...3.0.1#diff)
______________

### 3.0.0

__Changes__

- cli: add token argument to use instead of named auth
- cli: add rm-assets command
- cli: add rm-release command
- test: initialize tests
- [BREAK] most API: replaced username/password/OTP arguments
  by and instance of github.Client







__Contributors__

- mh-cbon

Released by mh-cbon, Fri 29 Jul 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/2.0.6...3.0.0#diff)
______________

### 2.0.6

__Changes__

- dl-assets: ensure only valid assets are scanned
- README: add repositories to install section

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 28 Jul 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/2.0.5...2.0.6#diff)
______________

### 2.0.5

__Changes__

- release: update release scripts
- travis: make use of an env variable for app name

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 28 Jul 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/2.0.4-beta...2.0.5#diff)
______________

### 2.0.4-beta

__Changes__

- release: update release scripts
- travis: make use of an env variable for app name

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 28 Jul 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/2.0.3...2.0.4-beta#diff)
______________

### 2.0.3

__Changes__

- release: update release scripts
- travis: set the GHTOKEN variable
- cli dl-assets: add skip-prerelease argument

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 28 Jul 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/2.0.2-beta4...2.0.3#diff)
______________

### 2.0.2-beta4

__Changes__

- travis: set the GHTOKEN variable

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 28 Jul 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/2.0.2-beta3...2.0.2-beta4#diff)
______________

### 2.0.2-beta3

__Changes__

- release: update release scripts

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 28 Jul 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/2.0.2-beta2...2.0.2-beta3#diff)
______________

### 2.0.2-beta2

__Changes__

- release: update release scripts

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 28 Jul 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/2.0.2-beta1...2.0.2-beta2#diff)
______________

### 2.0.2-beta1

__Changes__

- release: update release scripts

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 28 Jul 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/2.0.2-beta...2.0.2-beta1#diff)
______________

### 2.0.2-beta

__Changes__

- cli dl-assets: add skip-prerelease argument

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 28 Jul 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/2.0.1...2.0.2-beta#diff)
______________

### 2.0.1

__Changes__

- cli: add dl-assets command
- gh: Add methods to list public releases/assets and download asset

__Contributors__

- mh-cbon

Released by mh-cbon, Wed 27 Jul 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/2.0.0...2.0.1#diff)
______________

### 2.0.0

__Changes__

- [break] gh.CreateRelease: add new body argument to sets the body content of the release
- create-release: add changelog,c argument
- release: add changelog command to the creation of the release
- stringexec: add dependency

__Contributors__

- mh-cbon

Released by mh-cbon, Wed 27 Jul 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/1.0.2...2.0.0#diff)
______________

### 1.0.2

__Changes__

- github: update to latest github.com/google/go-github
- glide: fix missing semver dependency
- appveyor: fix indentation

__Contributors__

- mh-cbon

Released by mh-cbon, Fri 22 Jul 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/1.0.1...1.0.2#diff)
______________

### 1.0.1

__Changes__

- empty: trigger ci build

__Contributors__

- mh-cbon

Released by mh-cbon, Fri 22 Jul 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/1.0.0...1.0.1#diff)
______________

### 1.0.0

__Changes__

- create-release command: change --draft arg to take a string value (yes|no)
- glide: fix package name
- release: update release script

__Contributors__

- mh-cbon

Released by mh-cbon, Fri 22 Jul 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/0.0.3...1.0.0#diff)
______________

### 0.0.3

__Changes__

- Improve release scripts
- pass build
- Improve assets ulpoading

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 16 Jun 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/0.0.2...0.0.3#diff)
______________

### 0.0.2

__Changes__

- add publish script
- fix wrong imports
- README
- add release commands

__Contributors__

- mh-cbon

Released by mh-cbon, Wed 15 Jun 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/0.0.1...0.0.2#diff)
______________

### 0.0.1

__Changes__

- Initial release

__Contributors__

- mh-cbon

Released by mh-cbon, Tue 14 Jun 2016 -
[see the diff](https://github.com/mh-cbon/gh-api-cli/compare/6b4908780f93b52178e4fba49dd20ad2ce308649...0.0.1#diff)
______________


