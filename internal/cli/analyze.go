package cli

import (
	"context"
	"fmt"
	"os"

	anlzr "github.com/Fauxmen4/deps-check/internal/analyzer"
	"github.com/Fauxmen4/deps-check/internal/fetcher"
	"github.com/Fauxmen4/deps-check/internal/format"
	"github.com/Fauxmen4/deps-check/internal/parser"
	"github.com/Fauxmen4/deps-check/internal/proxy"
	"github.com/spf13/cobra"
)

type analyzeOptions struct {
	format      string // plain text (table), json
	directOnly  bool
	branch      string
	path        string
	showAll     bool
	concurrency int
}

func newAnalyzeCmd() *cobra.Command {
	opts := &analyzeOptions{}

	cmd := &cobra.Command{
		Use:   "analyze <repository-url>",
		Short: "Analyze dependencies of a GitHub repository",
		Long: `Analyze Go module dependencies and check for available updates.

Supports GitHub repositories via HTTPS or SSH URLs.

Examples:
  depscheck analyze https://github.com/gin-gonic/gin
  depscheck analyze git@github.com:user/repo.git
  depscheck analyze https://github.com/user/monorepo/tree/main/services/api
  depscheck analyze https://github.com/user/repo --format json
  depscheck analyze https://github.com/user/repo --direct-only --all`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnalyze(args[0], opts)
		},
	}

	cmd.Flags().StringVarP(&opts.format, "format", "f", "table", "output format: table, json")
	cmd.Flags().BoolVarP(&opts.directOnly, "direct-only", "d", false, "skip indirect dependencies")
	cmd.Flags().StringVarP(&opts.branch, "branch", "b", "", "branch to analyze (default: repository default)")
	cmd.Flags().StringVarP(&opts.path, "path", "p", "", "path to go.mod within the repository")
	cmd.Flags().BoolVar(&opts.showAll, "all", false, "show all dependencies, including up-to-date ones")
	cmd.Flags().IntVar(&opts.concurrency, "concurrency", 10, "max number of concurrent proxy requests")

	return cmd
}

func runAnalyze(repoURL string, opts *analyzeOptions) error {
	// parse url
	repo, err := parser.Parse(repoURL)
	if err != nil {
		return fmt.Errorf("failed to parse repository URL: %w", err)
	}

	// CLI flags overwrites URL parts
	if opts.branch != "" {
		repo.Branch = opts.branch
	}
	if opts.path != "" {
		repo.Path = opts.path
	}

	// fetch go.mod
	f, err := fetcher.NewGitHubFetcher()
	if err != nil {
		return fmt.Errorf("failed to fetch go.mod from repository: %w", err)
	}
	goMod, err := f.FetchGoMod(context.TODO(), repo)
	if err != nil {
		return fmt.Errorf("failed to fetch go.mod from repository: %w", err)
	}

	// analyze
	provider := proxy.NewVersionProvider()
	analyzer := anlzr.New(provider)
	moduleInfo, report, err := analyzer.Analyze(goMod, anlzr.Options{
		DirectOnly:  opts.directOnly,
		Concurrency: opts.concurrency,
	})
	if err != nil {
		return fmt.Errorf("failed to analyze module dependencies: %w", err)
	}

	// filter
	if !opts.showAll {
		report = anlzr.FilterUpdatable(report)
	}
	if len(report) == 0 {
		fmt.Println("All dependencies are up to date.")
		return nil
	}

	// format
	formatter, err := format.New(format.Format(opts.format))
	if err != nil {
		return fmt.Errorf("invalid format: %w", err)
	}
	return formatter.Write(os.Stdout, moduleInfo, report)
}
