package format

import (
	"encoding/json"
	"io"

	"github.com/Fauxmen4/deps-check/internal/analyzer"
)

type JSONFormatter struct{}

// jsonItem is a DependencyReport representation.
type jsonItem struct {
	Module   string              `json:"module"`
	Current  string              `json:"current"`
	Latest   string              `json:"latest,omitempty"`
	Update   analyzer.UpdateType `json:"update"`
	Indirect bool                `json:"indirect"`
	Error    string              `json:"error,omitempty"`
}

// Format writes report as JSON to w.
func (f *JSONFormatter) Write(w io.Writer, report []analyzer.DependencyReport) error {
	items := make([]jsonItem, len(report))
	for i, r := range report {
		items[i] = jsonItem{
			Module:   r.Module,
			Current:  r.Current,
			Latest:   r.Latest,
			Update:   r.Update,
			Indirect: r.Indirect,
		}
		if r.Err != nil {
			items[i].Error = r.Err.Error()
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(items)
}
