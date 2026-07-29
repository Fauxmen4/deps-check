package cli

import "github.com/spf13/cobra"

// Execute runs the root command.
func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:     "depscheck",
		Short:   "Analyze Go module dependencies for available updates",
	}

	root.AddCommand(newAnalyzeCmd())
	return root
}
