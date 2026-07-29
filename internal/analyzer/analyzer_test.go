package analyzer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseGoMod(t *testing.T) {
	t.Parallel()

	const gomodFull = `module github.com/example/repo

go 1.22

require (
	github.com/stretchr/testify v1.9.0
	github.com/spf13/cobra v1.8.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	golang.org/x/sys v0.15.0 // indirect
)
`

	const gomodEmpty = `module github.com/example/repo

go 1.22
`

	const gomodMinimal = `module github.com/example/repo
`

	const gomodInvalid = `this is not a valid go.mod`

	tests := []struct {
		name       string
		gomod      string
		directOnly bool
		wantInfo   ModuleInfo
		wantDeps   []Dependency
		wantErr    bool
	}{
		{
			name:       "all dependencies",
			gomod:      gomodFull,
			directOnly: false,
			wantInfo: ModuleInfo{
				Name:    "github.com/example/repo",
				Version: "1.22",
			},
			wantDeps: []Dependency{
				{Module: "github.com/stretchr/testify", Current: "v1.9.0", Indirect: false},
				{Module: "github.com/spf13/cobra", Current: "v1.8.0", Indirect: false},
				{Module: "github.com/davecgh/go-spew", Current: "v1.1.1", Indirect: true},
				{Module: "golang.org/x/sys", Current: "v0.15.0", Indirect: true},
			},
		},
		{
			name:       "direct only",
			gomod:      gomodFull,
			directOnly: true,
			wantInfo: ModuleInfo{
				Name:    "github.com/example/repo",
				Version: "1.22",
			},
			wantDeps: []Dependency{
				{Module: "github.com/stretchr/testify", Current: "v1.9.0", Indirect: false},
				{Module: "github.com/spf13/cobra", Current: "v1.8.0", Indirect: false},
			},
		},
		{
			name:       "no dependencies",
			gomod:      gomodEmpty,
			directOnly: false,
			wantInfo: ModuleInfo{
				Name:    "github.com/example/repo",
				Version: "1.22",
			},
			wantDeps: []Dependency{},
		},
		{
			name:       "no go directive",
			gomod:      gomodMinimal,
			directOnly: false,
			wantInfo: ModuleInfo{
				Name:    "github.com/example/repo",
				Version: "",
			},
			wantDeps: []Dependency{},
		},
		{
			name:    "invalid go.mod",
			gomod:   gomodInvalid,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotInfo, gotDeps, err := parseGoMod([]byte(tt.gomod), tt.directOnly)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseGoMod() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			if diff := cmp.Diff(tt.wantInfo, gotInfo); diff != "" {
				t.Errorf("ModuleInfo mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantDeps, gotDeps); diff != "" {
				t.Errorf("Dependencies mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestClassifyUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		current string
		latest  string
		want    UpdateType
	}{
		{
			name:    "same version",
			current: "v1.2.3",
			latest:  "v1.2.3",
			want:    UpdateNone,
		},
		{
			name:    "current is newer",
			current: "v1.3.0",
			latest:  "v1.2.5",
			want:    UpdateNone,
		},
		{
			name:    "patch update",
			current: "v1.2.3",
			latest:  "v1.2.5",
			want:    UpdatePatch,
		},
		{
			name:    "minor update",
			current: "v1.2.3",
			latest:  "v1.5.0",
			want:    UpdateMinor,
		},
		{
			name:    "major update",
			current: "v1.2.3",
			latest:  "v2.0.0",
			want:    UpdateMajor,
		},
		{
			name:    "v0 minor update",
			current: "v0.3.0",
			latest:  "v0.5.0",
			want:    UpdateMinor,
		},
		{
			name:    "v0 patch update",
			current: "v0.3.0",
			latest:  "v0.3.1",
			want:    UpdatePatch,
		},
		{
			name:    "invalid semver current",
			current: "not-a-version",
			latest:  "v1.0.0",
			want:    UpdateUnknown,
		},
		{
			name:    "invalid semver latest",
			current: "v1.0.0",
			latest:  "garbage",
			want:    UpdateUnknown,
		},
		{
			name:    "prerelease to stable",
			current: "v1.0.0-rc.1",
			latest:  "v1.0.0",
			want:    UpdatePatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := classifyUpdate(tt.current, tt.latest)
			if got != tt.want {
				t.Errorf("classifyUpdate(%q, %q) = %v, want %v",
					tt.current, tt.latest, got, tt.want)
			}
		})
	}
}
