package analyzer

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Fauxmen4/deps-check/internal/proxy"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

// UpdateType classifies the difference between a module's current and latest version.
type UpdateType string

const (
	UpdateNone    UpdateType = "none"
	UpdatePatch   UpdateType = "patch"
	UpdateMinor   UpdateType = "minor"
	UpdateMajor   UpdateType = "major"
	UpdateUnknown UpdateType = "unknown"
)

// VersionProvider resolves the latest version of a Go module.
type VersionProvider interface {
	LatestModuleVersion(moduleName string) (string, error)
}

// Dependency describes a single module extracted from go.mod.
type Dependency struct {
	Module   string
	Current  string
	Indirect bool
}

// DependencyReport is the analysis result for one module.
type DependencyReport struct {
	Module   string
	Current  string
	Latest   string
	Update   UpdateType
	Indirect bool
	Err      error
}

var (
	DefaultConcurrency = 10
)

// Options controls which deps are included and how they are checked.
type Options struct {
	DirectOnly  bool
	Concurrency int
}

// Analyzer compares dependencies from go.mod against latest versions returned by a VersionProvider.
type Analyzer struct {
	provider VersionProvider
}

func New(provider VersionProvider) *Analyzer {
	return &Analyzer{provider: provider}
}

// Analyze parses go.mod and returns update reports for its dependencies.
func (a *Analyzer) Analyze(goMod []byte, opts Options) ([]DependencyReport, error) {
	deps, err := parseDeps(goMod, opts.DirectOnly)
	if err != nil {
		return nil, err
	}

	concurrency := opts.Concurrency
	if concurrency <= 0 {
		concurrency = DefaultConcurrency
	}

	reports := make([]DependencyReport, len(deps))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i, dep := range deps {
		wg.Add(1)
		go func(i int, dep Dependency) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			reports[i] = a.analyze(dep)
		}(i, dep)
	}

	wg.Wait()
	return reports, nil
}

// analyze resolves the latest version for a single dependency.
func (a *Analyzer) analyze(dep Dependency) DependencyReport {
	report := DependencyReport{
		Module:   dep.Module,
		Current:  dep.Current,
		Indirect: dep.Indirect,
	}

	latest, err := a.provider.LatestModuleVersion(dep.Module)
	if err != nil {
		report.Update = UpdateUnknown

		if errors.Is(err, proxy.ErrNotFound) || errors.Is(err, proxy.ErrRetracted) {
			report.Err = err
			return report
		}
		report.Err = err
		return report
	}

	report.Latest = latest
	report.Update = classifyUpdate(dep.Current, latest)
	return report
}

// parseDeps extracts modules from go.mod, optionally skipping indirect ones.
func parseDeps(goMod []byte, directOnly bool) ([]Dependency, error) {
	f, err := modfile.Parse("go.mod", goMod, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse go.mod: %w", err)
	}

	deps := make([]Dependency, 0, len(f.Require))
	for _, r := range f.Require {
		if directOnly && r.Indirect {
			continue
		}
		deps = append(deps, Dependency{
			Module:   r.Mod.Path,
			Current:  r.Mod.Version,
			Indirect: r.Indirect,
		})
	}

	return deps, nil
}

// classifyUpdate returns the difference level between current and latest versions.
func classifyUpdate(current, latest string) UpdateType {
	if !semver.IsValid(current) || !semver.IsValid(latest) {
		return UpdateUnknown
	}

	if semver.Compare(current, latest) >= 0 {
		return UpdateNone
	}

	switch {
	case semver.Major(current) != semver.Major(latest):
		return UpdateMajor
	case semver.MajorMinor(current) != semver.MajorMinor(latest):
		return UpdateMinor
	default:
		return UpdatePatch
	}
}
