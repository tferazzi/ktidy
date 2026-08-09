package cmd

import (
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

// dryRunSkip prints a colored dry-run message and returns true when dry-run is active.
func dryRunSkip(cmd *cobra.Command, format string, args ...any) bool {
	if !dryRun {
		return false
	}
	colorDryRun.Fprintf(cmd.ErrOrStderr(), "[dry-run] "+format+"\n", args...)
	return true
}

// warnf prints a yellow warning line to stderr.
func warnf(cmd *cobra.Command, format string, args ...any) {
	colorWarning.Fprintf(cmd.ErrOrStderr(), "warning: "+format+"\n", args...)
}

// successf prints a green confirmation line to stdout.
func successf(cmd *cobra.Command, format string, args ...any) {
	colorSuccess.Fprintf(cmd.OutOrStdout(), format+"\n", args...)
}
