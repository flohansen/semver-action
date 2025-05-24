package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/oauth2"

	"github.com/flohansen/semver/internal/github"
)

var (
	ghRepo   = os.Getenv("GITHUB_REPOSITORY")
	ghSha    = os.Getenv("GITHUB_SHA")
	ghOutput = os.Getenv("GITHUB_OUTPUT")
)

type Config struct {
	Token string
}

type ActionApp struct {
	cfg Config
}

func NewAction(cfg Config) *ActionApp {
	return &ActionApp{
		cfg: cfg,
	}
}

func (a *ActionApp) Run(ctx context.Context) error {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: a.cfg.Token})
	tc := &http.Client{
		Transport: &oauth2.Transport{
			Source: ts,
			Base: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}

	repo, err := github.NewRepository(tc, ghRepo)
	if err != nil {
		return fmt.Errorf("could not create repository: %w", err)
	}

	commit, err := repo.GetLatestCommit(ctx, ghSha)
	if err != nil {
		return fmt.Errorf("could not parse latest commit: %w", err)
	}

	currentVersion, err := repo.GetLatestVersion(ctx)
	if err != nil {
		return fmt.Errorf("could not get latest version: %w", err)
	}

	fmt.Printf("Read current version: \033[32m%s\033[0m\n", currentVersion)
	fmt.Printf("Determine new version based on commit: \033[32m%s\033[0m\n", commit)

	newVersion := currentVersion
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

	fmt.Printf("Outputs:\n")
	fmt.Printf("  new-release: \033[32m%v\033[0m\n", newVersion != currentVersion)
	fmt.Printf("  new-release-version: \033[32m%s\033[0m\n", newVersion)

	if _, err := f.WriteString(fmt.Sprintf("new-release=%v\n", newVersion != currentVersion)); err != nil {
		return fmt.Errorf("could not write to GITHUB_OUTPUT: %w", err)
	}
	if _, err := f.WriteString(fmt.Sprintf("new-release-version=%s\n", newVersion)); err != nil {
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
