package gh

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// ClientFromCredentials creates a new github.Client instance,
// using a BasicAuthTransport.
func ClientFromCredentials(username string, password string, OTP string) *github.Client {

	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
		OTP:      strings.TrimSpace(OTP),
	}

	return github.NewClient(tp.Client())
}

// ClientFromToken creates a new github.Client instance,
// using an OAuth2 transport.
func ClientFromToken(token string) *github.Client {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	return github.NewClient(tc)
}

// AnonClient creates a new github.Client instance,
// using an anonymous transport.
func AnonClient() *github.Client {
	return github.NewClient(nil)
}

// Add creates a new authorization token on the remote.
func Add(client *github.Client, name string, permissions []string) (*github.Authorization, error) {

	authReq, notFound := GeneratePersonalAuthTokenRequest(name, permissions)
	if len(notFound) > 0 {
		return nil, errors.New("Unknown permissions: " + strings.Join(notFound, ","))
	}

	createdAuth, _, err := client.Authorizations.Create(authReq)

	return createdAuth, err
}

// List retrieve existing authorizations generated by gh-api-cli.
func List(client *github.Client) (map[string]*github.Authorization, error) {

	ret := make(map[string]*github.Authorization)

	opt := &github.ListOptions{Page: 1, PerPage: 200}
	got, _, err := client.Authorizations.List(opt)

	namedAuth := regexp.MustCompile(`^([^:]+): generated via (gh-auth|gh-api-cli)`)
	for _, v := range got {
		note := ""
		if v.Note != nil {
			note = *v.Note
		}
		if namedAuth.MatchString(note) {
			parts := namedAuth.FindAllStringSubmatch(note, -1)
			name := parts[0][1]
			ret[name] = v
		}
	}

	return ret, err
}

// Delete removes a personal access token on the remote
func Delete(client *github.Client, id int) error {

	_, err := client.Authorizations.Delete(id)

	return err
}

// GeneratePersonalAuthTokenRequest forges a request to create the personal access token on the remote
func GeneratePersonalAuthTokenRequest(name string, permissions []string) (*github.AuthorizationRequest, []string) {
	scopes := map[string]github.Scope{
		"user":             github.ScopeUser,
		"user:email":       github.ScopeUserEmail,
		"user:follow":      github.ScopeUserFollow,
		"public_repo":      github.ScopePublicRepo,
		"repo":             github.ScopeRepo,
		"repo_deployment":  github.ScopeRepoDeployment,
		"repo:status":      github.ScopeRepoStatus,
		"delete_repo":      github.ScopeDeleteRepo,
		"notifications":    github.ScopeNotifications,
		"gist":             github.ScopeGist,
		"read:repo_hook":   github.ScopeReadRepoHook,
		"write:repo_hook":  github.ScopeWriteRepoHook,
		"admin:repo_hook":  github.ScopeAdminRepoHook,
		"admin:org_hook":   github.ScopeAdminOrgHook,
		"read:org":         github.ScopeReadOrg,
		"write:org":        github.ScopeWriteOrg,
		"admin:org":        github.ScopeAdminOrg,
		"read:public_key":  github.ScopeReadPublicKey,
		"write:public_key": github.ScopeWritePublicKey,
		"admin:public_key": github.ScopeAdminPublicKey,
		"read:gpg_key":     github.ScopeReadGPGKey,
		"write:gpg_key":    github.ScopeWriteGPGKey,
		"admin:gpg_key":    github.ScopeAdminGPGKey,
	}
	notFound := make([]string, 0)
	auth := github.AuthorizationRequest{
		Note:        github.String(name + ": generated via gh-api-cli"),
		Scopes:      []github.Scope{},
		Fingerprint: github.String(name + time.Now().String()),
	}
	for _, p := range permissions {
		if val, ok := scopes[p]; ok {
			auth.Scopes = append(auth.Scopes, val)
		} else {
			notFound = append(notFound, p)
		}
	}
	return &auth, notFound
}

// ListReleases fetches first 200 releases on the remote.
func ListReleases(client *github.Client, owner string, repo string) ([]*github.RepositoryRelease, error) {

	opt := &github.ListOptions{Page: 1, PerPage: 200}
	got, _, err := client.Repositories.ListReleases(owner, repo, opt)

	return got, err
}

// GetReleaseByID gets a release by its id on the the remote
func GetReleaseByID(client *github.Client, owner string, repo string, id int) (*github.RepositoryRelease, error) {
	var ret *github.RepositoryRelease

	releases, err := ListReleases(client, owner, repo)
	if err != nil {
		return ret, err
	}
	for _, r := range releases {
		if *r.ID == id {
			ret = r
			break
		}
	}

	return ret, nil
}

// ListReleaseAssets retrieves the first 200 release assets on the remote
func ListReleaseAssets(client *github.Client, owner string, repo string, release github.RepositoryRelease) ([]*github.ReleaseAsset, error) {

	opt := &github.ListOptions{Page: 1, PerPage: 200}
	got, _, err := client.Repositories.ListReleaseAssets(owner, repo, *release.ID, opt)

	return got, err
}

// ReleaseExists tells if a release exits by its version on the remote
func ReleaseExists(client *github.Client, owner string, repo string, version string, draft bool) (bool, error) {

	releases, err := ListReleases(client, owner, repo)
	if err != nil {
		return true, err
	}

	exists := false
	for _, r := range releases {
		if (*r.TagName == version || *r.Name == version) && *r.Draft == draft {
			exists = true
			break
		}
	}

	return exists, err
}

// ReleaseID get ID of a release givne its version.
func ReleaseID(client *github.Client, owner string, repo string, version string) (int, error) {

	id := -1

	releases, err := ListReleases(client, owner, repo)
	if err != nil {
		return id, err
	}

	for _, r := range releases {
		if *r.TagName == version || *r.Name == version {
			id = *r.ID
			break
		}
	}

	if id == -1 {
		err = errors.New("Release '" + version + "' not found!")
	}

	return id, err
}

// CreateRelease creates a new release on the remote repository.
func CreateRelease(client *github.Client, owner, repo, version, author string, draft bool, body string) (*github.RepositoryRelease, error) {

	exists, err := ReleaseExists(client, owner, repo, version, draft)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("Release '" + version + "' already exists!")
	}
	v, err := semver.NewVersion(version)
	if err != nil {
		return nil, err
	}

	userAuthor := author
	if userAuthor == "" {
		userAuthor = owner
	}

	user, resp, err := client.Users.Get(userAuthor)
	if err != nil {
		return nil, fmt.Errorf("err: %v\nresp: %v", err, resp)
	}

	opt := &github.RepositoryRelease{
		Name:       github.String(version),
		TagName:    github.String(version),
		Draft:      github.Bool(draft),
		Body:       github.String(body),
		Prerelease: github.Bool(v.Prerelease() != ""),
		Author:     user,
	}
	release, _, err := client.Repositories.CreateRelease(owner, repo, opt)

	return release, err
}

// DeleteRelease deletes a release from the remote repository.
func DeleteRelease(client *github.Client, owner string, repo string, version string) error {

	id, err := ReleaseID(client, owner, repo, version)
	if err != nil {
		return err
	}
	_, err = client.Repositories.DeleteRelease(owner, repo, id)

	return err
}

// DeleteReleaseAsset deletes a release asset from the remote repository.
func DeleteReleaseAsset(client *github.Client, owner string, repo string, id int) error {

	_, err := client.Repositories.DeleteReleaseAsset(owner, repo, id)

	return err
}

// UploadReleaseAssets uploads multiple assets to the remote release identified by its version.
func UploadReleaseAssets(client *github.Client, owner string, repo string, version string, files []string) []error {

	errs := make([]error, 0)

	id, err := ReleaseID(client, owner, repo, version)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	c := make(chan error)
	UploadMultipleReleaseAssets(client, owner, repo, id, files, c)

	index := 0
	for upErr := range c {
		if upErr != nil {
			errs = append(errs, upErr)
		}
		index++
		if index == len(files)-1 {
			close(c)
		}
	}

	return errs
}

// UploadMultipleReleaseAssets uploads multiple assets to the remote release identified by its release ID.
func UploadMultipleReleaseAssets(client *github.Client, owner string, repo string, releaseID int, files []string, info chan<- error) {
	for index, file := range files {
		go func(index int, file string) {
			info <- UploadReleaseAsset(client, owner, repo, releaseID, file)
		}(index, file)
	}
}

// UploadReleaseAsset uploads one asset to the remote release identified by its release ID.
func UploadReleaseAsset(client *github.Client, owner string, repo string, releaseID int, file string) error {

	f, err := os.Open(file)
	defer f.Close()
	if err == nil {
		opt := &github.UploadOptions{Name: filepath.Base(file)}
		_, _, err = client.Repositories.UploadReleaseAsset(owner, repo, releaseID, opt, f)
	}

	return err
}

// DownloadAsset downloads an asset from a release, it handles redirect.
func DownloadAsset(url string, out io.Writer) error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	s3URL, err := res.Location()
	if err == nil && s3URL.String() != "" {
		res, err = http.Get(s3URL.String())
		if err != nil {
			return err
		}
		defer res.Body.Close()
	}

	_, err = io.Copy(out, res.Body)

	return err
}
