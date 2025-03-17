package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/flohansen/semver/internal/github"
)

var (
	ghRepo   = os.Getenv("GITHUB_REPOSITORY")
	ghSha    = os.Getenv("GITHUB_SHA")
	ghOutput = os.Getenv("GITHUB_OUTPUT")
)

type ActionFlags struct {
	OutputName string
}

type ActionApp struct {
	flags ActionFlags
}

func NewAction(flags ActionFlags) *ActionApp {
	return &ActionApp{
		flags: flags,
	}
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

	commit, err := repo.GetLatestCommit(ctx, ghSha)
	if err != nil {
		return fmt.Errorf("could not parse latest commit: %w", err)
	}

	version, err := repo.GetLatestVersion(ctx)
	if err != nil {
		return fmt.Errorf("could not get latest version: %w", err)
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
		return fmt.Errorf("could not open GITHUB_OUTPUT: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("%s=%s\n", a.flags.OutputName, version)); err != nil {
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
