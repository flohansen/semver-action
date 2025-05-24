package app

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/flohansen/semver/internal/domain"
	"github.com/flohansen/semver/internal/github"
)

var (
	ghRepo      = os.Getenv("GITHUB_REPOSITORY")
	ghSha       = os.Getenv("GITHUB_SHA")
	ghOutput    = os.Getenv("GITHUB_OUTPUT")
	ghEventPath = os.Getenv("GITHUB_EVENT_PATH")
)

type ActionApp struct {
}

func NewAction() *ActionApp {
	return &ActionApp{}
}

func getLatestCommit() (domain.Commit, error) {
	f, err := os.Open(ghEventPath)
	if err != nil {
		return domain.Commit{}, err
	}
	defer f.Close()

	var event github.Event
	if err := json.NewDecoder(f).Decode(&event); err != nil {
		return domain.Commit{}, err
	}

	return domain.NewCommitFromString(event.HeadCommit.Message)
}

func (a *ActionApp) Run(ctx context.Context) error {
	repo, err := github.NewRepository(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}, ghRepo)
	if err != nil {
		return fmt.Errorf("could not create repository: %w", err)
	}

	commit, err := getLatestCommit()
	if err != nil {
		return fmt.Errorf("could not parse latest commit: %w", err)
	}

	currentVersion, err := repo.GetLatestVersion(ctx)
	if err != nil {
		return fmt.Errorf("could not get latest version: %w", err)
	}
	newVersion := domain.Version{
		Major: currentVersion.Major,
		Minor: currentVersion.Minor,
		Patch: currentVersion.Patch,
	}

	fmt.Printf("Read current version: \033[32m%s\033[0m\n", currentVersion)
	fmt.Printf("Determine new version based on commit: \033[32m%s\033[0m\n", commit)

	if commit.IsBreaking {
		newVersion.IncMajor()
	} else if commit.Type == "feat" {
		newVersion.IncMinor()
	} else if commit.Type == "fix" {
		newVersion.IncPatch()
	}

	fmt.Printf("New version: \033[32m%s\033[0m\n", newVersion)

	f, err := os.OpenFile(ghOutput, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open GITHUB_OUTPUT: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("new-release=%v\n", newVersion == *currentVersion)); err != nil {
		return fmt.Errorf("could not write to GITHUB_OUTPUT: %w", err)
	}
	if _, err := f.WriteString(fmt.Sprintf("new-release-version=%s\n", currentVersion)); err != nil {
		return fmt.Errorf("could not write to GITHUB_OUTPUT: %w", err)
	}

	return nil
}

func SignalContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sig
		cancel()

		<-sig
		os.Exit(1)
	}()

	return ctx
}
