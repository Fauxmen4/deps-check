package parser

import "testing"

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    Repo
		wantErr bool
	}{
		{
			name:  "https simple",
			input: "https://github.com/gin-gonic/gin",
			want:  Repo{Owner: "gin-gonic", Name: "gin"},
		},
		{
			name:  "https with .git",
			input: "https://github.com/gin-gonic/gin.git",
			want:  Repo{Owner: "gin-gonic", Name: "gin"},
		},
		{
			name:  "https with branch",
			input: "https://github.com/user/repo/tree/develop",
			want:  Repo{Owner: "user", Name: "repo", Branch: "develop"},
		},
		{
			name:  "https with branch and path",
			input: "https://github.com/kubernetes/kubernetes/tree/master/staging/src/k8s.io/api",
			want: Repo{
				Owner:  "kubernetes",
				Name:   "kubernetes",
				Branch: "master",
				Path:   "staging/src/k8s.io/api",
			},
		},
		{
			name:  "ssh",
			input: "git@github.com:gin-gonic/gin.git",
			want:  Repo{Owner: "gin-gonic", Name: "gin"},
		},
		{
			name:    "unsupported host",
			input:   "https://gitlab.com/user/repo",
			wantErr: true,
		},
		{
			name:    "missing repo",
			input:   "https://github.com/user",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("Parse() = %+v, want %+v", got, tt.want)
			}
		})
	}
}