package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/flohansen/semver/internal/github"
)

var (
	ghRepo   = os.Getenv("GITHUB_REPOSITORY")
	ghSha    = os.Getenv("GITHUB_SHA")
	ghOutput = os.Getenv("GITHUB_OUTPUT")
)

func main() {
	ctx := context.Background()

	repo, err := github.NewRepository(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}, ghRepo)
	if err != nil {
		fmt.Printf("could not create repository: %s", err)
		os.Exit(1)
	}

	commit, err := repo.GetLatestCommit(ctx, ghSha)
	if err != nil {
		fmt.Printf("could not parse latest commit: %s", err)
		os.Exit(1)
	}

	version, err := repo.GetLatestVersion(ctx)
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

	f, err := os.OpenFile(ghOutput, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
