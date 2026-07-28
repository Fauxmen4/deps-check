package parser

import (
	"fmt"
	"net/url"
	"strings"
)

// Repo is a repository metadata representation
type Repo struct {
	Owner  string
	Name   string
	Branch string
	Path   string
}

// Parse parses repository url 
func Parse(rawURL string) (Repo, error) {
	if strings.HasPrefix(rawURL, "git@") {
		return parseSSH(rawURL)
	}
	if strings.HasPrefix(rawURL, "https://") {
		return parseHTTPS(rawURL)
	}

	return Repo{}, fmt.Errorf("unsupported URL scheme: %s", rawURL)
}

func parseHTTPS(rawURL string) (Repo, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return Repo{}, fmt.Errorf("invalid URL :%w", err)
	}

	if u.Host != "github.com" {
		return Repo{}, fmt.Errorf("only github.com is supported, got: %s", u.Host)
	}

	pathParts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(pathParts) < 2 {
		return Repo{}, fmt.Errorf("URL must contain owner and repository: %s", rawURL)
	}

	repo := Repo{
		Owner: pathParts[0],
		Name:  strings.TrimSuffix(pathParts[1], ".git"),
	}

	if len(pathParts) >= 4 && pathParts[2] == "tree" {
		repo.Branch = pathParts[3]
		if len(pathParts) > 4 {
			repo.Path = strings.Join(pathParts[4:], "/")
		}
	}

	return repo, nil
}

func parseSSH(rawURL string) (Repo, error) {
	const prefix = "git@github.com:"

	if !strings.HasPrefix(rawURL, prefix) {
		return Repo{}, fmt.Errorf("only github.com is supported: %s", rawURL)
	}

	path := strings.TrimPrefix(rawURL, prefix)
	path = strings.TrimSuffix(path, ".git")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(pathParts) != 2 {
		return Repo{}, fmt.Errorf("SSH URL must be in format git@github.com:owner/repo.git")
	}

	return Repo{
		Owner: pathParts[0],
		Name:  pathParts[1],
	}, nil
}
