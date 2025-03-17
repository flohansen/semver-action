package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/v39/github"
)

var (
	ghOwner = os.Getenv("GITHUB_REPOSITORY_OWNER")
	ghRepo  = os.Getenv("GITHUB_REPOSITORY")
	ghSha   = os.Getenv("GITHUB_SHA")
)

func main() {
	ctx := context.Background()
	client := github.NewClient(&http.Client{})

	commit, err := getLatestCommit(ctx, client)
	if err != nil {
		fmt.Printf("could not parse latest commit: %s", err)
		os.Exit(1)
	}

	version, err := getLatestVersion(ctx, client)
	if err != nil {
		fmt.Printf("could not get latest version: %s", err)
		os.Exit(1)
	}

	fmt.Printf("Read current version: \033[32m%s\033[0m\n", version)
	fmt.Printf("Determine new version based on commit: \033[32m%s\033[0m\n", commit)

	if commit.IsBreaking {
		version.IncMajor()
	} else if commit.Type == "feat" {
		version.IncMinor()
	} else if commit.Type == "fix" {
		version.IncPatch()
	}

	fmt.Printf("New version: \033[32m%s\033[0m\n", version)

	f, err := os.OpenFile(os.Getenv("GITHUB_OUTPUT"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("could not open GITHUB_OUTPUT: %s", err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("new_version=%s\n", version)); err != nil {
		fmt.Printf("could not write to GITHUB_OUTPUT: %s", err)
		os.Exit(1)
	}
}

var convCommReg = regexp.MustCompile(`([a-zA-Z ]+)(!)?(\(.*\))?:(.*)`)

type Commit struct {
	Raw        string
	Type       string
	Scope      string
	Message    string
	IsBreaking bool
}

func (c *Commit) String() string {
	return c.Raw
}

func NewCommitFromString(str string) (*Commit, error) {
	match := convCommReg.FindStringSubmatch(str)
	if len(match) == 0 {
		return nil, errors.New("invalid commit message")
	}

	c := &Commit{}
	c.Raw = str
	c.Type = strings.ToLower(match[1])
	c.IsBreaking = c.Type == "breaking change" || match[2] == "!"
	return c, nil
}

func getLatestCommit(ctx context.Context, client *github.Client) (*Commit, error) {
	c, _, err := client.Repositories.GetCommit(ctx, ghOwner, ghRepo, ghSha, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting commit: %w", err)
	}

	commit, err := NewCommitFromString(strings.TrimSpace(*c.Commit.Message))
	if err != nil {
		return nil, fmt.Errorf("error reading commit: %w", err)
	}

	return commit, nil
}

var verReg = regexp.MustCompile(`(v)?([0-9]+)\.([0-9]+)\.([0-9]+)`)

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v *Version) IncMajor() {
	v.Major++
	v.Minor = 0
	v.Patch = 0
}

func (v *Version) IncMinor() {
	v.Minor++
	v.Patch = 0
}

func (v *Version) IncPatch() {
	v.Patch++
}

func (v *Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func NewVersionFromString(str string) (*Version, error) {
	match := verReg.FindStringSubmatch(str)
	if len(match) == 0 {
		return nil, errors.New("invalid version string")
	}

	major, err := strconv.Atoi(match[2])
	if err != nil {
		return nil, fmt.Errorf("could not parse major version: %w", err)
	}

	minor, err := strconv.Atoi(match[3])
	if err != nil {
		return nil, fmt.Errorf("could not parse minot version: %w", err)
	}

	patch, err := strconv.Atoi(match[4])
	if err != nil {
		return nil, fmt.Errorf("could not parse patch version: %w", err)
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

func getLatestVersion(ctx context.Context, client *github.Client) (*Version, error) {
	tags, _, err := client.Repositories.ListTags(ctx, ghOwner, ghRepo, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting tags: %w", err)
	}

	tag := "v0.0.0"
	if len(tags) > 0 {
		tag = *tags[0].Name
	}

	version, err := NewVersionFromString(tag)
	if err != nil {
		return nil, fmt.Errorf("error reading version: %w", err)
	}

	return version, nil
}
