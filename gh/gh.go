package gh

import (
	"time"
	"regexp"
	"strings"
	"errors"

	"github.com/google/go-github/github"
)

func Add (username string, password string, OTP string, name string, permissions []string) (*github.Authorization, error) {

  tp := github.BasicAuthTransport{
    Username: strings.TrimSpace(username),
    Password: strings.TrimSpace(password),
    OTP:      strings.TrimSpace(OTP),
  }

  client := github.NewClient(tp.Client())

  authReq, notFound := GeneratePersonalAuthTokenRequest(name, permissions)
  if len(notFound)>0 {
    return nil, errors.New("Unknown permissions: " + strings.Join(notFound, ","))
  }

  createdAuth, _, err := client.Authorizations.Create(authReq)

  return createdAuth, err
}

func List (username string, password string, OTP string) (map[string]*github.Authorization, error) {

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
    if v.Note!=nil {
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

func Delete (username string, password string, OTP string, id int) error {

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
    "user":               github.ScopeUser,
    "user:email":         github.ScopeUserEmail,
    "user:follow":        github.ScopeUserFollow,
    "public_repo":        github.ScopePublicRepo,
    "repo":               github.ScopeRepo,
    "repo_deployment":    github.ScopeRepoDeployment,
    "repo:status":        github.ScopeRepoStatus,
    "delete_repo":        github.ScopeDeleteRepo,
    "notifications":      github.ScopeNotifications,
    "gist":               github.ScopeGist,
    "read:repo_hook":     github.ScopeReadRepoHook,
    "write:repo_hook":    github.ScopeWriteRepoHook,
    "admin:repo_hook":    github.ScopeAdminRepoHook,
    "admin:org_hook":     github.ScopeAdminOrgHook,
    "read:org":           github.ScopeReadOrg,
    "write:org":          github.ScopeWriteOrg,
    "admin:org":          github.ScopeAdminOrg,
    "read:public_key":    github.ScopeReadPublicKey,
    "write:public_key":   github.ScopeWritePublicKey,
    "admin:public_key":   github.ScopeAdminPublicKey,
    "read:gpg_key":       github.ScopeReadGPGKey,
    "write:gpg_key":      github.ScopeWriteGPGKey,
    "admin:gpg_key":      github.ScopeAdminGPGKey,
  }
  notFound := make([]string, 0)
	auth := github.AuthorizationRequest{
		Note:        github.String(name+": generated via gh-auth"),
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
