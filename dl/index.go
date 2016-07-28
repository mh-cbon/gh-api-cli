package dl

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"github.com/mh-cbon/gh-api-cli/gh"
)

type Asset struct {
	SourceUrl  string
	TargetFile string
	Name       string
	Arch       string
	System     string
	Ext        string
	Version    string
}

// Select assets of given releases matching glob,
// forge out path and url for each asset
func SelectAssets(owner string, repo string, glob string, out string, releases []*github.RepositoryRelease) ([]*Asset, error) {
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
			if *a.State == "uploaded" && r.MatchString(*a.Name) {
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

// Select releases matching constraint
func SelectReleases(constraint string, skipPrerelease bool, releases []*github.RepositoryRelease) ([]*github.RepositoryRelease, error) {
	ret := make([]*github.RepositoryRelease, 0)
	if constraint == "latest" {
		release, _ := SelectLatestRelease(skipPrerelease, releases)
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
			if c.Check(v) && (!skipPrerelease || skipPrerelease && v.Prerelease() == "") {
				ret = append(ret, r)
			}
		}
	}
	return ret, nil
}

// Select release which are not prerelease
func SelectNonPrerelease(releases []*github.RepositoryRelease) ([]*github.RepositoryRelease, error) {
	ret := make([]*github.RepositoryRelease, 0)
	for _, r := range releases {
		v, err := semver.NewVersion(*r.TagName)
		if err != nil {
			continue
		}
		if v.Prerelease() == "" {
			ret = append(ret, r)
		}
	}
	return ret, nil
}

// Select only latest release according to semver sort
func SelectLatestRelease(skipPrerelease bool, releases []*github.RepositoryRelease) (*github.RepositoryRelease, error) {
	var release *github.RepositoryRelease
	for _, r := range releases {
		v, err := semver.NewVersion(*r.TagName)
		if err != nil {
			continue
		}
		if release == nil && (!skipPrerelease || skipPrerelease && v.Prerelease() == "") {
			release = r
			continue
		}
		if release != nil {
			v2, err := semver.NewVersion(*release.TagName)
			if err != nil {
				continue
			}
			if v.GreaterThan(v2) && (!skipPrerelease || skipPrerelease && v2.Prerelease() == "") {
				release = r
			}
		}
	}
	return release, nil
}

func DownloadAsset(asset *Asset) error {
	fmt.Println("Downloading " + asset.Name + " to " + asset.TargetFile + ", version=" + asset.Version)
	dir := filepath.Dir(asset.TargetFile)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return errors.New("Failed to create directory " + dir + "!\n" + err.Error())
	}
	f, err := os.Create(asset.TargetFile)
	if err != nil {
		return errors.New("Failed to open file " + asset.TargetFile + "!\n" + err.Error())
	}
	err = gh.DownloadAsset(asset.SourceUrl, f)
	if err != nil {
		return errors.New("Failed to download file!\n" + err.Error())
	}
	return nil
}
