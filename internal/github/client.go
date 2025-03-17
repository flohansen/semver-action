package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/flohansen/semver/internal/domain"
	"github.com/google/go-github/v39/github"
)

type Repository struct {
	client *github.Client
	owner  string
	repo   string
}

func NewRepository(client *http.Client, repo string) (*Repository, error) {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid repository string, expected format: '<owner>/<repository>'")
	}

	return &Repository{
		client: github.NewClient(client),
		owner:  parts[0],
		repo:   parts[1],
	}, nil
}

func (r *Repository) GetLatestVersion(ctx context.Context) (*domain.Version, error) {
	tags, _, err := r.client.Repositories.ListTags(ctx, r.owner, r.repo, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting tags: %w", err)
	}

	tag := "v0.0.0"
	if len(tags) > 0 {
		tag = *tags[0].Name
	}

	version, err := domain.NewVersionFromString(tag)
	if err != nil {
		return nil, fmt.Errorf("error reading version: %w", err)
	}

	return version, nil
}

func (r *Repository) GetLatestCommit(ctx context.Context, sha string) (*domain.Commit, error) {
	c, _, err := r.client.Repositories.GetCommit(ctx, r.owner, r.repo, sha, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting commit: %w", err)
	}

	commit, err := domain.NewCommitFromString(strings.TrimSpace(*c.Commit.Message))
	if err != nil {
		return nil, fmt.Errorf("error reading commit: %w", err)
	}

	return commit, nil
}
