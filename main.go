package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"github.com/mattn/go-zglob"
	"github.com/mh-cbon/gh-api-cli/gh"
	"github.com/mh-cbon/gh-api-cli/local"
	"github.com/mh-cbon/stringexec"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
)

var VERSION = "0.0.0"

func main() {

	app := cli.NewApp()
	app.Name = "gh-api-cli"
	app.Version = VERSION
	app.Usage = "Github api command line client"
	app.UsageText = "gh-api-cli <cmd> <options>"
	app.Commands = []cli.Command{
		{
			Name:   "add-auth",
			Usage:  "Add a new authorization",
			Action: add,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "username, u",
					Value: "",
					Usage: "Github username",
				},
				cli.StringFlag{
					Name:  "password, p",
					Value: "",
					Usage: "Github password",
				},
				cli.StringFlag{
					Name:  "name, n",
					Value: "",
					Usage: "Name of the authorization",
				},
				cli.StringSliceFlag{
					Name:  "rights, r",
					Value: &cli.StringSlice{},
					Usage: "Permissions to set",
				},
			},
		},
		{
			Name:   "list-auth",
			Usage:  "List authorizations",
			Action: list,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "username, u",
					Value: "",
					Usage: "Github username",
				},
				cli.StringFlag{
					Name:  "password, p",
					Value: "",
					Usage: "Github password",
				},
			},
		},
		{
			Name:   "rm-auth",
			Usage:  "Remove an authorization",
			Action: rm,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "username, u",
					Value: "",
					Usage: "Github username",
				},
				cli.StringFlag{
					Name:  "password, p",
					Value: "",
					Usage: "Github password",
				},
				cli.StringFlag{
					Name:  "name, n",
					Value: "",
					Usage: "Name of the authorization",
				},
			},
		},
		{
			Name:   "get-auth",
			Usage:  "Get token from a locally saved authorization",
			Action: get,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Value: "",
					Usage: "Name of the authorization",
				},
			},
		},
		{
			Name:   "create-release",
			Usage:  "Create a release",
			Action: createRelease,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Value: "",
					Usage: "Name of the authorization",
				},
				cli.StringFlag{
					Name:  "owner, o",
					Value: "",
					Usage: "Repo owner",
				},
				cli.StringFlag{
					Name:  "repository, r",
					Value: "",
					Usage: "Repo name",
				},
				cli.StringFlag{
					Name:  "ver",
					Value: "",
					Usage: "Version name",
				},
				cli.StringFlag{
					Name:  "author, a",
					Value: "",
					Usage: "Release author name",
				},
				cli.StringFlag{
					Name:  "email, e",
					Value: "",
					Usage: "Release author email",
				},
				cli.StringFlag{
					Name:  "draft, d",
					Value: "no",
					Usage: "Make a draft release",
				},
				cli.StringFlag{
					Name:  "changelog, c",
					Value: "",
					Usage: "A command to generate the description body of the release",
				},
			},
		},
		{
			Name:   "upload-release-asset",
			Usage:  "Upload assets to a release",
			Action: uploadReleaseAsset,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Value: "",
					Usage: "Name of the authorization",
				},
				cli.StringFlag{
					Name:  "glob, g",
					Value: "",
					Usage: "Glob pattern of files to upload",
				},
				cli.StringFlag{
					Name:  "owner, o",
					Value: "",
					Usage: "Repo owner",
				},
				cli.StringFlag{
					Name:  "repository, r",
					Value: "",
					Usage: "Repo name",
				},
				cli.StringFlag{
					Name:  "ver",
					Value: "",
					Usage: "Version name",
				},
			},
		},
		{
			Name:   "dl-assets",
			Usage:  "Download assets",
			Action: downloadAssets,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "glob, g",
					Value: "",
					Usage: "Glob pattern of files to download",
				},
				cli.StringFlag{
					Name:  "out",
					Value: "%f",
					Usage: "Out format to write files",
				},
				cli.StringFlag{
					Name:  "owner, o",
					Value: "",
					Usage: "Repo owner",
				},
				cli.StringFlag{
					Name:  "repository, r",
					Value: "",
					Usage: "Repo name",
				},
				cli.StringFlag{
					Name:  "ver",
					Value: "",
					Usage: "Version constraint",
				},
				cli.StringFlag{
					Name:  "skip-prerelease",
					Value: "no",
					Usage: "Skip prerelease releases (yes|no)",
				},
			},
		},
	}

	app.Run(os.Args)
}

func add(c *cli.Context) error {

	name := c.String("name")
	perms := make([]string, 0)
	for _, p := range c.StringSlice("rights") {
		perms = append(perms, string(p))
	}

	if len(name) == 0 {
		return cli.NewExitError("You must provide an authorization name", 1)
	}
	if len(perms) == 0 {
		return cli.NewExitError("You must provide permissions", 1)
	}

	username := getUsername(c.String("username"))
	password := getPassword(c.String("password"))

	auths, err := gh.List(username, password, "")
	if _, ok := err.(*github.TwoFactorAuthError); err != nil && ok {
		otp := queryOtp()
		auths, err = gh.List(username, password, otp)
		if err != nil {
			fmt.Println(err)
			return cli.NewExitError("Could not list current authorizations!", 1)
		}
	}

	//-auths
	if _, ok := auths[name]; ok {
		return cli.NewExitError("Authorization "+name+" already exists!", 1)
	}

	createdAuth, err := gh.Add(username, password, "", name, perms)
	if _, ok := err.(*github.TwoFactorAuthError); err != nil && ok {
		otp := queryOtp()
		createdAuth, err = gh.Add(username, password, otp, name, perms)
	}

	if err != nil {
		fmt.Println(err)
		return cli.NewExitError("The request could not be completed!", 1)
	}

	err = local.Add(name, createdAuth)
	if err != nil {
		fmt.Println(err)
		fmt.Println("The new token was not saved on your local!")
	}

	fmt.Println("New authorization '" + name + "' created!")
	fmt.Println(string(*createdAuth.Token))

	return nil
}

func list(c *cli.Context) error {

	username := getUsername(c.String("username"))
	password := getPassword(c.String("password"))

	auths, err := gh.List(username, password, "")
	if _, ok := err.(*github.TwoFactorAuthError); err != nil && ok {
		otp := queryOtp()
		auths, err = gh.List(username, password, otp)
	}
	if err != nil {
		return cli.NewExitError("Could not list current authorizations!", 1)
	}

	saved := local.Read()

	for name, auth := range auths {
		if val, ok := saved[name]; ok {
			if val.Token != nil {
				auth.Token = github.String(*val.Token)
			} else {
				auth.Token = github.String("Unknown on your local")
			}
		} else {
			auth.Token = github.String("Unknown on your local")
		}
	}

	jsonContent, err := jsonString(auths)
	if err != nil {
		fmt.Println(err)
		return cli.NewExitError("There was an error while printing your results!", 1)
	}
	fmt.Println(jsonContent)

	return nil
}

func rm(c *cli.Context) error {
	name := c.String("name")

	if len(name) == 0 {
		return cli.NewExitError("You must provide a name", 1)
	}

	otp := ""
	username := getUsername(c.String("username"))
	password := getPassword(c.String("password"))

	auths, err := gh.List(username, password, "")
	if _, ok := err.(*github.TwoFactorAuthError); err != nil && ok {
		otp = queryOtp()
		auths, err = gh.List(username, password, otp)
	}
	if err != nil {
		fmt.Println(err)
		return cli.NewExitError("Could not list current authorizations!", 1)
	}

	if val, ok := auths[name]; ok {
		err = gh.Delete(username, password, otp, *val.ID)
		if err != nil {
			fmt.Println(err)
			return cli.NewExitError("The deletion failed!", 1)
		}
	} else {
		return cli.NewExitError("Authorization "+name+" was not found!", 1)
	}

	err = local.Remove(name)
	if err != nil {
		fmt.Println(err)
		fmt.Println("The authorization was not removed on your local!")
	}

	fmt.Println("Deleted authorization: " + name)

	return nil
}

func get(c *cli.Context) error {
	name := c.String("name")

	if len(name) == 0 {
		return cli.NewExitError("You must provide a name", 1)
	}

	auth, err := local.Get(name)
	if err != nil {
		fmt.Println(err)
		return cli.NewExitError("The authorization '"+name+"' was not found on your local!", 1)
	}

	if auth.Token == nil {
		return cli.NewExitError("The authorization '"+name+"' does not have token!", 1)
	}
	fmt.Println(*auth.Token)

	return nil
}

func createRelease(c *cli.Context) error {
	name := c.String("name")
	owner := c.String("owner")
	repo := c.String("repository")
	ver := c.String("ver")
	author := c.String("author")
	email := c.String("email")
	draft := c.String("draft")
	changelog := c.String("changelog")
	isDraft := false
	body := ""

	if len(name) == 0 {
		return cli.NewExitError("You must provide an authorization name", 1)
	}
	if len(owner) == 0 {
		return cli.NewExitError("You must provide the repository owner", 1)
	}
	if len(repo) == 0 {
		return cli.NewExitError("You must provide a repository name", 1)
	}
	if len(ver) == 0 {
		return cli.NewExitError("You must provide a version", 1)
	}
	if len(author) != 0 && len(email) == 0 {
		return cli.NewExitError("You must provide an email", 1)
	}
	if len(author) == 0 && len(email) > 0 {
		return cli.NewExitError("You must provide an author", 1)
	}
	if draft == "yes" || draft == "1" || draft == "true" {
		isDraft = true
	}

	if changelog != "" {
		oCmd, err := stringexec.Command(changelog)
		if err != nil {
			fmt.Println(err)
			return cli.NewExitError("The changelog was not properly generated!", 1)
		}
		oCmd.Stderr = os.Stderr
		out, err := oCmd.Output()
		if err != nil {
			fmt.Println(err)
			return cli.NewExitError("The changelog was not properly generated!", 1)
		}
		body = string(out)
	}

	auth, err := local.Get(name)
	if err != nil {
		fmt.Println(err)
		return cli.NewExitError("The authorization '"+name+"' was not found on your local!", 1)
	}

	if auth.Token == nil {
		return cli.NewExitError("The authorization '"+name+"' does not have token!", 1)
	}

	release, err := gh.CreateRelease(*auth.Token, owner, repo, ver, author, email, isDraft, body)
	if err != nil {
		fmt.Println(err)
		return cli.NewExitError("The release was not created successfully!", 1)
	}

	jsonContent, _ := jsonString(release)
	fmt.Println(jsonContent)

	return nil
}

func uploadReleaseAsset(c *cli.Context) error {
	name := c.String("name")
	glob := c.String("glob")
	owner := c.String("owner")
	repo := c.String("repository")
	ver := c.String("ver")

	if len(name) == 0 {
		return cli.NewExitError("You must provide an authorization name", 1)
	}
	if len(glob) == 0 {
		return cli.NewExitError("You must provide a pattern to glob", 1)
	}
	if len(owner) == 0 {
		return cli.NewExitError("You must provide a repository owner", 1)
	}
	if len(repo) == 0 {
		return cli.NewExitError("You must provide a repository name", 1)
	}
	if len(ver) == 0 {
		return cli.NewExitError("You must provide a release version", 1)
	}

	auth, err := local.Get(name)
	if err != nil {
		fmt.Println(err)
		return cli.NewExitError("The authorization '"+name+"' was not found on your local!", 1)
	}

	if auth.Token == nil {
		return cli.NewExitError("The authorization '"+name+"' does not have token!", 1)
	}

	paths, err := zglob.Glob(glob)
	if len(paths) == 0 {
		return cli.NewExitError("Your glob pattern did not selected any files.", 1)
	}

	token := *auth.Token
	id, err := gh.ReleaseId(token, owner, repo, ver)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	errs := make([]error, 0)
	for _, file := range paths {
		fmt.Println("Uploading " + file)
		err := gh.UploadReleaseAsset(token, owner, repo, id, file)
		if err != nil {
			fmt.Println("Failed")
			errs = append(errs, err)
			fmt.Println(err)
		} else {
			fmt.Println("Done")
		}
	}

	if len(errs) > 0 {
		return cli.NewExitError("There were errors while uploading assets.", 1)
	} else {
		fmt.Println("Assets uploaded!")
	}

	return nil
}

func downloadAssets(c *cli.Context) error {
	glob := c.String("glob")
	out := c.String("out")
	owner := c.String("owner")
	repo := c.String("repository")
	ver := c.String("ver")
	sp := c.String("skip-prerelease")
  skipPrerelease := false

	if len(owner) == 0 {
		return cli.NewExitError("You must provide a repository owner", 1)
	}
	if len(repo) == 0 {
		return cli.NewExitError("You must provide a repository name", 1)
	}
  if sp=="yes" || sp=="true" || sp=="1" {
    skipPrerelease = true
  }

	releases, err := gh.ListPublicReleases(owner, repo)
	if err != nil {
		fmt.Println(err)
		return cli.NewExitError("could not list releases of this repository "+owner+"/"+repo+"!", 1)
	}

	if ver != "" {
		releases, err = selectReleases(ver, skipPrerelease, releases)
		if err != nil {
			fmt.Println(err)
			return cli.NewExitError("Failed to select release for this constraint "+ver+"!", 1)
		}
	} else if skipPrerelease {
		releases, err = selectNonPrerelease(releases)
		if err != nil {
			fmt.Println(err)
			return cli.NewExitError("Failed to select release for this constraint "+ver+"!", 1)
		}
  }

	if len(releases) == 0 {
		fmt.Println("No releases selected!")
		return nil
	}

	assets, err := selectAssets(owner, repo, glob, out, releases)
	if err != nil {
		fmt.Println(err)
		return cli.NewExitError("Failed to select assets for this glob "+glob+"!", 1)
	}

	if len(assets) == 0 {
		fmt.Println("No assets selected!")
		return nil
	}

	for _, a := range assets {
		fmt.Println("Downloading " + a.Name + " to " + a.TargetFile+", version="+a.Version)
		dir := filepath.Dir(a.TargetFile)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println(err)
			return cli.NewExitError("Failed to create directory "+dir+"!", 1)
		}
		f, err := os.Create(a.TargetFile)
		if err != nil {
			fmt.Println(err)
			return cli.NewExitError("Failed to open file "+a.TargetFile+"!", 1)
		}
		err = gh.DownloadAsset(a.SourceUrl, f)
		if err != nil {
			fmt.Println(err)
			return cli.NewExitError("Failed to download file!", 1)
		}
	}

	fmt.Println("All done!")

	return nil
}

type Asset struct {
	SourceUrl  string
	TargetFile string
	Name       string
	Arch       string
	System     string
	Ext        string
	Version    string
}

func selectAssets(owner string, repo string, glob string, out string, releases []*github.RepositoryRelease) ([]*Asset, error) {
	ret := make([]*Asset, 0)
	r, _ := regexp.Compile(".+")
	if glob != "" {
		var err error
		glob = strings.Replace(glob, "*", ".+", -1)
		r, err = regexp.Compile("(?i)^" + glob + "$")
		if err != nil {
			return ret, err
		}
	}
	for _, release := range releases {
		assets, err := gh.ListReleaseAssets(owner, repo, *release)
		if err != nil {
			return ret, err
		}
		for _, a := range assets {
			if r.MatchString(*a.Name) {
				asset := &Asset{}
				asset.Name = *a.Name
				asset.Version = *release.TagName
				asset.SourceUrl = *a.BrowserDownloadURL
				ret = append(ret, asset)
			}
		}
	}
	isWin, _ := regexp.Compile("(?i)win(dows)?[-_.]")
	isDarwin, _ := regexp.Compile("(?i)(darwin|mac)[-_.]")
	isLinux, _ := regexp.Compile("(?i)(linux|ubuntu|debian|fedora|arch|gentoo)[-_.]")
	isWinExt, _ := regexp.Compile("(?i).+[.](exe|msi)$")
	isMacExt, _ := regexp.Compile("(?i).+[.](dmg)$")
	isLinuxExt, _ := regexp.Compile("(?i).+[.](deb|rpm)$")
	is386, _ := regexp.Compile("(?i)[-_.](386|i386)[-_.]")
	isAmd64, _ := regexp.Compile("(?i)[-_.](amd64|x86_64)[-_.]")
	for _, a := range ret {
		a.Ext = filepath.Ext(a.Name)
		if isWin.MatchString(a.Name) {
			a.System = "windows"
		} else if isDarwin.MatchString(a.Name) {
			a.System = "darwin"
		} else if isLinux.MatchString(a.Name) {
			a.System = "linux"
		}
		if is386.MatchString(a.Name) {
			a.Arch = "386"
		} else if isAmd64.MatchString(a.Name) {
			a.Arch = "amd64"
		}
		if a.System == "" {
			if isWinExt.MatchString(a.Name) {
				a.System = "windows"
			} else if isMacExt.MatchString(a.Name) {
				a.System = "darwin"
			} else if isLinuxExt.MatchString(a.Name) {
				a.System = "linux"
			} else {
				a.System = "unknown"
			}
		}
	}
	for _, a := range ret {
		a.TargetFile = out
		/*
		   %f: full filename
		   %o: repository owner
		   %r: repository name
		   %e: file extension, minus dot prefix, detected JIT
		   %s: target system (windows, darwin, linux), detected JIT
		   %a: architecture (amd64, 386), detected JIT
		   %v: version the asset is attached to
		*/
		e := strings.TrimPrefix(a.Ext, ".")
		a.TargetFile = strings.Replace(a.TargetFile, "%f", a.Name, -1)
		a.TargetFile = strings.Replace(a.TargetFile, "%o", owner, -1)
		a.TargetFile = strings.Replace(a.TargetFile, "%r", repo, -1)
		a.TargetFile = strings.Replace(a.TargetFile, "%e", e, -1)
		a.TargetFile = strings.Replace(a.TargetFile, "%s", a.System, -1)
		a.TargetFile = strings.Replace(a.TargetFile, "%a", a.Arch, -1)
		a.TargetFile = strings.Replace(a.TargetFile, "%v", a.Version, -1)
	}
	return ret, nil
}

func selectReleases(constraint string, skipPrerelease bool, releases []*github.RepositoryRelease) ([]*github.RepositoryRelease, error) {
	ret := make([]*github.RepositoryRelease, 0)
	if constraint == "latest" {
		release, _ := selectLatestRelease(skipPrerelease, releases)
		if release != nil {
			ret = append(ret, release)
		}
	} else {
		c, err := semver.NewConstraint(constraint)
		if err != nil {
			return ret, err
		}
		for _, r := range releases {
			v, err := semver.NewVersion(*r.TagName)
			if err != nil {
				continue
			}
			if c.Check(v) && (!skipPrerelease || skipPrerelease && v.Prerelease()=="") {
				ret = append(ret, r)
			}
		}
	}
	return ret, nil
}

func selectNonPrerelease(releases []*github.RepositoryRelease) ([]*github.RepositoryRelease, error) {
	ret := make([]*github.RepositoryRelease, 0)
  for _, r := range releases {
    v, err := semver.NewVersion(*r.TagName)
    if err != nil {
      continue
    }
    if v.Prerelease()=="" {
      ret = append(ret, r)
    }
  }
	return ret, nil
}

func selectLatestRelease(skipPrerelease bool, releases []*github.RepositoryRelease) (*github.RepositoryRelease, error) {
	var release *github.RepositoryRelease
	for _, r := range releases {
		v, err := semver.NewVersion(*r.TagName)
		if err != nil {
			continue
		}
		if release == nil && (!skipPrerelease || skipPrerelease && v.Prerelease()=="") {
			release = r
			continue
		}
    if release != nil {
  		v2, err := semver.NewVersion(*release.TagName)
  		if err != nil {
  			continue
  		}
  		if v.GreaterThan(v2) && (!skipPrerelease || skipPrerelease && v2.Prerelease()=="") {
  			release = r
  		}
    }
	}
	return release, nil
}

func getUsername(username string) string {
	if strings.TrimSpace(username) == "" {
		username = queryUsername()
	}
	return username
}
func getPassword(password string) string {
	if strings.TrimSpace(password) == "" {
		password = queryPassword()
	}
	return password
}

func queryUsername() string {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("GitHub Username: ")
	username, _ := r.ReadString('\n')
	return strings.TrimSpace(username)
}

func queryOtp() string {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("\nGitHub One Time Password: ")
	otp, _ := r.ReadString('\n')
	return strings.TrimSpace(otp)
}

func queryPassword() string {
	fmt.Print("GitHub Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println("")
	return string(bytePassword)
}

func jsonString(some interface{}) (string, error) {
	jsonContent, err := json.MarshalIndent(some, "", "    ")
	return string(jsonContent), err
}

func exitWithError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
