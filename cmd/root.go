package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dryRun bool

var rootCmd = &cobra.Command{
	Use:   "ktidy",
	Short: "Merge, split, import, export and group kubeconfig files",
}

func Execute(version string) {
	rootCmd.Version = version
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "print what would happen without writing any files")
}

// dryRunSkip prints a dry-run message and returns true when dry-run is active.
// Usage: if dryRunSkip(cmd, "would write %s", path) { return nil }
func dryRunSkip(cmd *cobra.Command, format string, args ...any) bool {
	if !dryRun {
		return false
	}
	fmt.Fprintf(cmd.ErrOrStderr(), "[dry-run] "+format+"\n", args...)
	return true
}
