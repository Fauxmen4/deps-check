package fetcher

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Fauxmen4/deps-check/internal/parser"
	"github.com/google/go-github/v89/github"
)

var (
	ErrNotFound     = errors.New("go.mod not found")
	ErrRateLimited  = errors.New("GitHub API rate limit exceeded")
	ErrUnauthorized = errors.New("unauthorized")
)

// GitHubFetcher retrieves go.mod files from GitHub repositories via the GitHub REST API.
type GitHubFetcher struct {
	client *github.Client
}

func NewGitHubFetcher() (*GitHubFetcher, error) {
	client, err := github.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to init github client: %w", err)
	}
	return &GitHubFetcher{client: client}, nil
}

// FetchGoMod retrieves the contents of the go.mod file for the given repository. 
// If repo.Path is set, go.mod is looked up under that subdirectory; 
// otherwise, the repository root is used.
func (f *GitHubFetcher) FetchGoMod(ctx context.Context, repo parser.Repo) ([]byte, error) {
	path := repo.Path
	if path == "" {
		path = "go.mod"
	} else {
		path = path + "/go.mod"
	}

	var opts *github.RepositoryContentGetOptions
	if repo.Branch != "" {
		opts = &github.RepositoryContentGetOptions{Ref: repo.Branch}
	}

	file, _, resp, err := f.client.Repositories.GetContents(
		ctx, repo.Owner, repo.Name, path, opts,
	)
	if err != nil {
		return nil, mapError(err, resp, path)
	}

	if file == nil {
		return nil, fmt.Errorf("%w: path %q is a directory, not a file", ErrNotFound, path)
	}

	content, err := file.GetContent()
	if err != nil {
		return nil, fmt.Errorf("decoding go.mod content: %w", err)
	}

	return []byte(content), nil
}

func mapError(err error, resp *github.Response, path string) error {
	if resp == nil {
		return fmt.Errorf("fetching go.mod: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		return fmt.Errorf("%w at %q: %w", ErrNotFound, path, err)
	case http.StatusUnauthorized, http.StatusForbidden:
		if _, ok := err.(*github.RateLimitError); ok {
			return fmt.Errorf("%w: %w", ErrRateLimited, err)
		}
		return fmt.Errorf("%w: %w", ErrUnauthorized, err)
	default:
		return fmt.Errorf("fetching go.mod: %w", err)
	}
}
