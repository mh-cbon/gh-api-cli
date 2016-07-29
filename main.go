package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"syscall"

	"github.com/google/go-github/github"
	"github.com/mattn/go-zglob"
	"github.com/mh-cbon/gh-api-cli/dl"
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
					Usage: "Name of the authorization to create",
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
					Usage: "Name of the authorization to delete",
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
					Usage: "Name of the authorization to use for identification",
				},
				cli.StringFlag{
					Name:  "token, t",
					Value: "",
					Usage: "Value of a personal access token",
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
			Name:   "rm-release",
			Usage:  "Delete a release",
			Action: rmRelease,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Value: "",
					Usage: "Name of the authorization to use for identification",
				},
				cli.StringFlag{
					Name:  "token, t",
					Value: "",
					Usage: "Value of a personal access token",
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
			Name:   "upload-release-asset",
			Usage:  "Upload assets to a release",
			Action: uploadReleaseAsset,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Value: "",
					Usage: "Name of the authorization to use for identification",
				},
				cli.StringFlag{
					Name:  "token, t",
					Value: "",
					Usage: "Value of a personal access token",
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
					Name:  "name, n",
					Value: "",
					Usage: "Name of the authorization to use for identification",
				},
				cli.StringFlag{
					Name:  "token, t",
					Value: "",
					Usage: "Value of a personal access token",
				},
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
		{
			Name:   "rm-assets",
			Usage:  "Delete assets",
			Action: rmAssets,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Value: "",
					Usage: "Name of the authorization to use for identification",
				},
				cli.StringFlag{
					Name:  "token, t",
					Value: "",
					Usage: "Value of a personal access token",
				},
				cli.StringFlag{
					Name:  "glob, g",
					Value: "",
					Usage: "Glob pattern of files to download",
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
			},
		},
	}

	app.Run(os.Args)
}

// Create a new authorization on remote save it locally
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
	client := gh.ClientFromCredentials(username, password, "")

	auths, err := gh.List(client)
	if _, ok := err.(*github.TwoFactorAuthError); err != nil && ok {
		otp := queryOtp()
		client = gh.ClientFromCredentials(username, password, otp)
		auths, err = gh.List(client)
		if err != nil {
			fmt.Println(err)
			return cli.NewExitError("Could not list current authorizations!", 1)
		}
	}

	//-auths
	if _, ok := auths[name]; ok {
		return cli.NewExitError("Authorization "+name+" already exists!", 1)
	}

	createdAuth, err := gh.Add(client, name, perms)
	if _, ok := err.(*github.TwoFactorAuthError); err != nil && ok {
		otp := queryOtp()
		client = gh.ClientFromCredentials(username, password, otp)
		createdAuth, err = gh.Add(client, name, perms)
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

// List authorizations from remote and add token values saved on local
func list(c *cli.Context) error {

	var client *github.Client
	username := getUsername(c.String("username"))
	password := getPassword(c.String("password"))
	client = gh.ClientFromCredentials(username, password, "")

	auths, err := gh.List(client)
	if _, ok := err.(*github.TwoFactorAuthError); err != nil && ok {
		otp := queryOtp()
		client = gh.ClientFromCredentials(username, password, otp)
		auths, err = gh.List(client)
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

// Delete an authorization from local and remote
func rm(c *cli.Context) error {
	name := c.String("name")

	if len(name) == 0 {
		return cli.NewExitError("You must provide a name of autorization to delete", 1)
	}

	var client *github.Client
	username := getUsername(c.String("username"))
	password := getPassword(c.String("password"))
	client = gh.ClientFromCredentials(username, password, "")

	auths, err := gh.List(client)
	if _, ok := err.(*github.TwoFactorAuthError); err != nil && ok {
		otp := queryOtp()
		client = gh.ClientFromCredentials(username, password, otp)
		auths, err = gh.List(client)
	}
	if err != nil {
		fmt.Println(err)
		return cli.NewExitError("Could not list current authorizations!", 1)
	}

	if val, ok := auths[name]; ok {
		err = gh.Delete(client, *val.ID)
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

// Print a token from an authorization saved on local
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

// Create a gh release on remote
func createRelease(c *cli.Context) error {
	name := c.String("name")
	token := c.String("token")
	owner := c.String("owner")
	repo := c.String("repository")
	ver := c.String("ver")
	author := c.String("author")
	email := c.String("email")
	draft := c.String("draft")
	changelog := c.String("changelog")
	isDraft := false
	body := ""

	if len(name)+len(token) == 0 {
		return cli.NewExitError("You must provide an authorization (--name or --token)", 1)
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

	tokenAuth := token
	var client *github.Client
	if len(name) > 0 {
		auth, err := local.Get(name)
		if err != nil {
			fmt.Println(err)
			return cli.NewExitError("The authorization '"+name+"' was not found on your local!", 1)
		}
		if auth.Token == nil {
			return cli.NewExitError("The authorization '"+name+"' does not have token!", 1)
		}
		tokenAuth = *auth.Token
	}

	client = gh.ClientFromToken(tokenAuth)

	release, err := gh.CreateRelease(client, owner, repo, ver, author, email, isDraft, body)
	if err != nil {
		fmt.Println(err)
		return cli.NewExitError("The release was not created successfully!", 1)
	}

	jsonContent, _ := jsonString(release)
	fmt.Println(jsonContent)

	return nil
}

// Delete a gh release on remote
func rmRelease(c *cli.Context) error {
	name := c.String("name")
	token := c.String("token")
	owner := c.String("owner")
	repo := c.String("repository")
	ver := c.String("ver")

	if len(name)+len(token) == 0 {
		return cli.NewExitError("You must provide an authorization (--name or --token)", 1)
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

	tokenAuth := token
	var client *github.Client
	if len(name) > 0 {
		auth, err := local.Get(name)
		if err != nil {
			fmt.Println(err)
			return cli.NewExitError("The authorization '"+name+"' was not found on your local!", 1)
		}
		if auth.Token == nil {
			return cli.NewExitError("The authorization '"+name+"' does not have token!", 1)
		}
		tokenAuth = *auth.Token
	}

	client = gh.ClientFromToken(tokenAuth)

	err := gh.DeleteRelease(client, owner, repo, ver)
	if err != nil {
		fmt.Println(err)
		return cli.NewExitError("The release was not deleted successfully!", 1)
	}

	fmt.Println("Release deleted with success!")

	return nil
}

// Upload asset to a release
func uploadReleaseAsset(c *cli.Context) error {
	name := c.String("name")
	token := c.String("token")
	glob := c.String("glob")
	owner := c.String("owner")
	repo := c.String("repository")
	ver := c.String("ver")

	if len(name)+len(token) == 0 {
		return cli.NewExitError("You must provide an authorization (--name or --token)", 1)
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

	tokenAuth := token
	var client *github.Client
	if len(name) > 0 {
		auth, err := local.Get(name)
		if err != nil {
			fmt.Println(err)
			return cli.NewExitError("The authorization '"+name+"' was not found on your local!", 1)
		}
		if auth.Token == nil {
			return cli.NewExitError("The authorization '"+name+"' does not have token!", 1)
		}
		tokenAuth = *auth.Token
	}

	client = gh.ClientFromToken(tokenAuth)

	paths, err := zglob.Glob(glob)
	if len(paths) == 0 {
		return cli.NewExitError("Your glob pattern did not selected any files.", 1)
	}

	id, err := gh.ReleaseId(client, owner, repo, ver)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	errs := make([]error, 0)
	for _, file := range paths {
		fmt.Println("Uploading " + file)
		err := gh.UploadReleaseAsset(client, owner, repo, id, file)
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

// Upload asset to a release
func rmAssets(c *cli.Context) error {
	name := c.String("name")
	token := c.String("token")
	glob := c.String("glob")
	owner := c.String("owner")
	repo := c.String("repository")
	ver := c.String("ver")

	if len(name)+len(token) == 0 {
		return cli.NewExitError("You must provide an authorization (--name or --token)", 1)
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
	if len(glob) == 0 {
		glob = "*"
	}

	tokenAuth := token
	var client *github.Client
	if len(name) > 0 {
		auth, err := local.Get(name)
		if err != nil {
			fmt.Println(err)
			return cli.NewExitError("The authorization '"+name+"' was not found on your local!", 1)
		}
		if auth.Token == nil {
			return cli.NewExitError("The authorization '"+name+"' does not have token!", 1)
		}
		tokenAuth = *auth.Token
	}

	client = gh.ClientFromToken(tokenAuth)

	id, err := gh.ReleaseId(client, owner, repo, ver)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	release, err := gh.GetReleaseById(client, owner, repo, id)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	if release == nil {
		return cli.NewExitError("release '"+ver+"' was not found!", 1)
	}

	r, _ := regexp.Compile(".+")
	if glob != "" {
		var err error
		glob = strings.Replace(glob, "*", ".+", -1)
		r, err = regexp.Compile("(?i)^" + glob + "$")
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}
	assets, err := gh.ListReleaseAssets(client, owner, repo, *release)
	for _, a := range assets {
		if r.MatchString(*a.Name) {
			if err = gh.DeleteReleaseAsset(client, owner, repo, *a.ID); err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			fmt.Println("Removed '" + (*a.Name) + "'")
		}
	}

	fmt.Println("All done!")

	return nil
}

// Download asset from a release
func downloadAssets(c *cli.Context) error {
	name := c.String("name")
	token := c.String("token")
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
	if sp == "yes" || sp == "true" || sp == "1" {
		skipPrerelease = true
	}

	tokenAuth := token
	var client *github.Client
	if len(name) > 0 {
		auth, err := local.Get(name)
		if err != nil {
			fmt.Println(err)
			return cli.NewExitError("The authorization '"+name+"' was not found on your local!", 1)
		}
		if auth.Token == nil {
			return cli.NewExitError("The authorization '"+name+"' does not have token!", 1)
		}
		tokenAuth = *auth.Token
	}

	if len(tokenAuth) > 0 {
		client = gh.ClientFromToken(tokenAuth)
	} else {
		client = gh.AnonClient()
	}

	releases, err := gh.ListPublicReleases(client, owner, repo)
	if err != nil {
		fmt.Println(err)
		return cli.NewExitError("could not list releases of this repository "+owner+"/"+repo+"!", 1)
	}

	if ver != "" {
		releases, err = dl.SelectReleases(ver, skipPrerelease, releases)
		if err != nil {
			fmt.Println(err)
			return cli.NewExitError("Failed to select release for this constraint "+ver+"!", 1)
		}
	} else if skipPrerelease {
		releases, err = dl.SelectNonPrerelease(releases)
		if err != nil {
			fmt.Println(err)
			return cli.NewExitError("Failed to select release for this constraint "+ver+"!", 1)
		}
	}

	if len(releases) == 0 {
		fmt.Println("No releases selected!")
		return nil
	}

	assets, err := dl.SelectAssets(client, owner, repo, glob, out, releases)
	if err != nil {
		fmt.Println(err)
		return cli.NewExitError("Failed to select assets for this glob "+glob+"!", 1)
	}

	if len(assets) == 0 {
		fmt.Println("No assets selected!")
		return nil
	}

	for _, a := range assets {
		fmt.Println("Downloading " + a.Name + " to " + a.TargetFile + ", version=" + a.Version)
		err := dl.DownloadAsset(a)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	fmt.Println("All done!")

	return nil
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
