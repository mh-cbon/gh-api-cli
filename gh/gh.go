package gh

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"io"
  "net/http"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Add a new personal access tokens
func Add(username string, password string, OTP string, name string, permissions []string) (*github.Authorization, error) {

	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
		OTP:      strings.TrimSpace(OTP),
	}

	client := github.NewClient(tp.Client())

	authReq, notFound := GeneratePersonalAuthTokenRequest(name, permissions)
	if len(notFound) > 0 {
		return nil, errors.New("Unknown permissions: " + strings.Join(notFound, ","))
	}

	createdAuth, _, err := client.Authorizations.Create(authReq)

	return createdAuth, err
}

// List personal access tokens generated via gh-api-cli on the remote
func List(username string, password string, OTP string) (map[string]*github.Authorization, error) {

	ret := make(map[string]*github.Authorization)

	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
		OTP:      strings.TrimSpace(OTP),
	}

	client := github.NewClient(tp.Client())

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

// Delete a personal access token on the remote
func Delete(username string, password string, OTP string, id int) error {

	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
		OTP:      strings.TrimSpace(OTP),
	}

	client := github.NewClient(tp.Client())

	_, err := client.Authorizations.Delete(id)

	return err
}

// Forge a request to create the personal access token on the remote
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

// List all releases on the remote
func ListReleases(token string, owner string, repo string) ([]*github.RepositoryRelease, error) {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	opt := &github.ListOptions{Page: 1, PerPage: 200}
	got, _, err := client.Repositories.ListReleases(owner, repo, opt)

	return got, err
}

// List public releases on the remote
func ListPublicReleases(owner string, repo string) ([]*github.RepositoryRelease, error) {

	client := github.NewClient(nil)

	opt := &github.ListOptions{Page: 1, PerPage: 200}
	got, _, err := client.Repositories.ListReleases(owner, repo, opt)

	return got, err
}

// List public release assets on the remote
func ListReleaseAssets(owner string, repo string, release github.RepositoryRelease) ([]*github.ReleaseAsset, error) {

	client := github.NewClient(nil)

	opt := &github.ListOptions{Page: 1, PerPage: 200}
	got, _, err := client.Repositories.ListReleaseAssets(owner, repo, *release.ID, opt)

	return got, err
}

// Tells if a release exitst on the remote
func ReleaseExists(token string, owner string, repo string, version string, draft bool) (bool, error) {

	releases, err := ListReleases(token, owner, repo)
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

// Transform a version string into its id
func ReleaseId(token string, owner string, repo string, version string) (int, error) {

	id := -1

	releases, err := ListReleases(token, owner, repo)
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

// Create a new release on the remote
func CreateRelease(token string, owner string, repo string, version string, authorName string, authorEmail string, draft bool, body string) (*github.RepositoryRelease, error) {

	exists, err := ReleaseExists(token, owner, repo, version, draft)
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

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryRelease{
		Name:       github.String(version),
		TagName:    github.String(version),
		Draft:      github.Bool(draft),
		Body:       github.String(body),
		Prerelease: github.Bool(v.Prerelease() != ""),
		Author: &github.CommitAuthor{
			Name:  github.String(authorName),
			Email: github.String(authorEmail),
		},
	}
	release, _, err := client.Repositories.CreateRelease(owner, repo, opt)

	return release, err
}

// Upload multiple assets the remote release
func UploadReleaseAssets(token string, owner string, repo string, version string, files []string) []error {

	errs := make([]error, 0)

	id, err := ReleaseId(token, owner, repo, version)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	c := make(chan error)
	UploadMultipleReleaseAssets(token, owner, repo, id, files, c)

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
func UploadMultipleReleaseAssets(token string, owner string, repo string, releaseId int, files []string, info chan<- error) {
	for index, file := range files {
		go func(index int, file string) {
			info <- UploadReleaseAsset(token, owner, repo, releaseId, file)
		}(index, file)
	}
}

// Upload one asset on the remote
func UploadReleaseAsset(token string, owner string, repo string, releaseId int, file string) error {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	f, err := os.Open(file)
	defer f.Close()
	if err == nil {
		opt := &github.UploadOptions{Name: filepath.Base(file)}
		_, _, err = client.Repositories.UploadReleaseAsset(owner, repo, releaseId, opt, f)
	}

	return err
}

// Download an asset from a release, handles redirect.
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
  if err==nil && s3URL.String()!="" {
    res, err = http.Get(s3URL.String())
    if err != nil {
        return err
    }
    defer res.Body.Close()
  }

  _, err = io.Copy(out, res.Body)

  return err
}
