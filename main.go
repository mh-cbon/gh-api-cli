package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/google/go-github/github"
	"github.com/mattn/go-zglob"
	"github.com/mh-cbon/gh-api-cli/GenVersionFile"
	"github.com/mh-cbon/gh-api-cli/gh"
	"github.com/mh-cbon/gh-api-cli/local"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {

	app := cli.NewApp()
	app.Name = "gh-api-cli"
	app.Version = GenVersionFile.Version()
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
				cli.BoolFlag{
					Name:  "draft, d",
					Usage: "Make a draft release",
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
	draft := c.Bool("draft")

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

	auth, err := local.Get(name)
	if err != nil {
		fmt.Println(err)
		return cli.NewExitError("The authorization '"+name+"' was not found on your local!", 1)
	}

	if auth.Token == nil {
		return cli.NewExitError("The authorization '"+name+"' does not have token!", 1)
	}

	release, err := gh.CreateRelease(*auth.Token, owner, repo, ver, author, email, draft)
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

	id, err := ReleaseId(token, owner, repo, ver)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

  errs := make([]error, 0)
  token := *auth.Token
  for index, file := range paths {
    fmt.Println("Uploading "+file)
    err := UploadReleaseAsset(token, owner, repo, id, file)
    if err!= nil {
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
