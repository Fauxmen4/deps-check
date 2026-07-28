package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const defaultProxyURL = "https://proxy.golang.org"

var (
	ErrNotFound = errors.New("module not found")

	ErrRetracted = errors.New("module retracted")
)

// VersionProvider resolves the latest version of a Go module using the
// module proxy protocol described at https://proxy.golang.org.
type VersionProvider struct {
	baseURL string
	client  *http.Client
}

// NewVersionProvider returns a VersionProvider that talks to the default
// public Go module proxy (proxy.golang.org).
func NewVersionProvider() *VersionProvider {
	return &VersionProvider{
		baseURL: defaultProxyURL,
	}
}

type versionInfo struct {
	Version string `json:"Version"`
}

// LatestModuleVersion returns the latest known version of the given module.
// It returns ErrNotFound if the proxy has no record of the module, and
// ErrGone if the module has been retracted.
func (vp *VersionProvider) LatestModuleVersion(module string) (string, error) {
	url := fmt.Sprintf("%s/%s/@latest", vp.baseURL, module)

	resp, err := vp.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch latest version of %s: %w", module, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return "", fmt.Errorf("%w: %s", ErrNotFound, module)
	case http.StatusGone:
		return "", fmt.Errorf("%w: %s", ErrRetracted, module)
	default:
		return "", fmt.Errorf("unexpected status %d for %s", resp.StatusCode, module)
	}

	var v versionInfo
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return "", fmt.Errorf("failed to decode response for %s: %w", module, err)
	}
	if v.Version == "" {
		return "", fmt.Errorf("empty version in response for %s", module)
	}

	return v.Version, nil
}
