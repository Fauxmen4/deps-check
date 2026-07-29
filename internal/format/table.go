package format

import (
	"fmt"
	"io"

	"github.com/Fauxmen4/deps-check/internal/analyzer"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type TableFormatter struct{}

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Padding(0, 1)
	cellStyle   = lipgloss.NewStyle().Padding(0, 1)
	borderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	updateColors = map[analyzer.UpdateType]lipgloss.Color{
		analyzer.UpdatePatch:   lipgloss.Color("2"), // green
		analyzer.UpdateMinor:   lipgloss.Color("3"), // yellow
		analyzer.UpdateMajor:   lipgloss.Color("1"), // red
		analyzer.UpdateUnknown: lipgloss.Color("8"), // gray
	}
)

// Format writes reports as a styled table to w.
func (f *TableFormatter) Write(w io.Writer, moduleInfo analyzer.ModuleInfo, reports []analyzer.DependencyReport) error {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(borderStyle).
		Headers("MODULE", "CURRENT", "LATEST", "UPDATE", "INDIRECT").
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			// Colorize the UPDATE column (index 3).
			if col == 3 && row < len(reports) {
				if c, ok := updateColors[reports[row].Update]; ok {
					return cellStyle.Foreground(c)
				}
			}
			return cellStyle
		})

	for _, r := range reports {
		latest := r.Latest
		update := string(r.Update)
		if r.Err != nil {
			latest = "-"
			update = fmt.Sprintf("error: %s", r.Err)
		}
		indirect := "no"
		if r.Indirect {
			indirect = "yes"
		}
		t.Row(r.Module, r.Current, latest, update, indirect)
	}

	_, err := fmt.Fprintf(w, "Module: %s, version: %s\nDependencies:\n", moduleInfo.Name, moduleInfo.Version)
	_ = err
	_, err = fmt.Fprintln(w, t.Render())
	return err
}
