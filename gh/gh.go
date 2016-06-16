package gh

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

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

	namedAuth := regexp.MustCompile(`^([^:]+): generated via gh-auth`)
	for _, v := range got {
		note := ""
		if v.Note != nil {
			note = *v.Note
		}
		if namedAuth.MatchString(note) {
			parts := namedAuth.FindAllStringSubmatch(note, -1)
			name := parts[0][1]
			ret[name] = &v
		}
	}

	return ret, err
}

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

//
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
		Note:        github.String(name + ": generated via gh-auth"),
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

func ListReleases(token string, owner string, repo string) ([]github.RepositoryRelease, error) {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	opt := &github.ListOptions{Page: 1, PerPage: 200}
	got, _, err := client.Repositories.ListReleases(owner, repo, opt)

	return got, err
}

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

func CreateRelease(token string, owner string, repo string, version string, authorName string, authorEmail string, draft bool) (*github.RepositoryRelease, error) {

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
		Prerelease: github.Bool(v.Prerelease() != ""),
		Author: &github.CommitAuthor{
			Name:  github.String(authorName),
			Email: github.String(authorEmail),
		},
	}
	release, _, err := client.Repositories.CreateRelease(owner, repo, opt)

	return release, err
}

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
