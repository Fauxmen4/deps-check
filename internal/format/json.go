package format

import (
	"encoding/json"
	"io"

	"github.com/Fauxmen4/deps-check/internal/analyzer"
)

type JSONFormatter struct{}

type jsonReport struct {
	Name         string       `json:"name"`
	Version      string       `json:"version"`
	Dependencies []dependency `json:"dependencies"`
}

type dependency struct {
	Module   string              `json:"module"`
	Current  string              `json:"current"`
	Latest   string              `json:"latest,omitempty"`
	Update   analyzer.UpdateType `json:"update"`
	Indirect bool                `json:"indirect"`
	Error    string              `json:"error,omitempty"`
}

// Format writes report as JSON to w.
func (f *JSONFormatter) Write(w io.Writer, moduleInfo analyzer.ModuleInfo, report []analyzer.DependencyReport) error {
	deps := make([]dependency, len(report))
	for i, r := range report {
		deps[i] = dependency{
			Module:   r.Module,
			Current:  r.Current,
			Latest:   r.Latest,
			Update:   r.Update,
			Indirect: r.Indirect,
		}
		if r.Err != nil {
			deps[i].Error = r.Err.Error()
		}
	}

	r := jsonReport{
		Name: moduleInfo.Name, 
		Version: moduleInfo.Version, 
		Dependencies: deps,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
