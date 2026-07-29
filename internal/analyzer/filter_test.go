package analyzer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestFilterUpdatable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []DependencyReport
		want  []DependencyReport
	}{
		{
			name:  "empty input",
			input: []DependencyReport{},
			want:  []DependencyReport{},
		},
		{
			name: "all up to date",
			input: []DependencyReport{
				{Module: "a", Update: UpdateNone},
				{Module: "b", Update: UpdateNone},
			},
			want: []DependencyReport{},
		},
		{
			name: "all need updates",
			input: []DependencyReport{
				{Module: "a", Update: UpdatePatch},
				{Module: "b", Update: UpdateMinor},
				{Module: "c", Update: UpdateMajor},
			},
			want: []DependencyReport{
				{Module: "a", Update: UpdatePatch},
				{Module: "b", Update: UpdateMinor},
				{Module: "c", Update: UpdateMajor},
			},
		},
		{
			name: "mixed",
			input: []DependencyReport{
				{Module: "a", Update: UpdateNone},
				{Module: "b", Update: UpdatePatch},
				{Module: "c", Update: UpdateNone},
				{Module: "d", Update: UpdateMajor},
			},
			want: []DependencyReport{
				{Module: "b", Update: UpdatePatch},
				{Module: "d", Update: UpdateMajor},
			},
		},
		{
			name: "unknown kept",
			input: []DependencyReport{
				{Module: "a", Update: UpdateNone},
				{Module: "b", Update: UpdateUnknown},
			},
			want: []DependencyReport{
				{Module: "b", Update: UpdateUnknown},
			},
		},
		{
			name: "preserves order",
			input: []DependencyReport{
				{Module: "z", Update: UpdatePatch},
				{Module: "a", Update: UpdateNone},
				{Module: "m", Update: UpdateMinor},
			},
			want: []DependencyReport{
				{Module: "z", Update: UpdatePatch},
				{Module: "m", Update: UpdateMinor},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := FilterUpdatable(tt.input)
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("FilterUpdatable() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
