package analyzer

import (
	"testing"
)

func TestParseDeps(t *testing.T) {
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

	const gomodInvalid = `this is not a valid go.mod`

	tests := []struct {
		name       string
		gomod      string
		directOnly bool
		want       []Dependency
		wantErr    bool
	}{
		{
			name:       "all dependencies",
			gomod:      gomodFull,
			directOnly: false,
			want: []Dependency{
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
			want: []Dependency{
				{Module: "github.com/stretchr/testify", Current: "v1.9.0", Indirect: false},
				{Module: "github.com/spf13/cobra", Current: "v1.8.0", Indirect: false},
			},
		},
		{
			name:       "no dependencies",
			gomod:      gomodEmpty,
			directOnly: false,
			want:       []Dependency{},
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

			got, err := parseDeps([]byte(tt.gomod), tt.directOnly)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseDeps() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			if len(got) != len(tt.want) {
				t.Fatalf("got %d deps, want %d\ngot:  %+v\nwant: %+v",
					len(got), len(tt.want), got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("dep[%d] = %+v, want %+v", i, got[i], tt.want[i])
				}
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