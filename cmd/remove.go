package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tferazzi/ktidy/internal/kubeconfig"
)

var removeCmd = &cobra.Command{
	Use:   "remove [flags] CONTEXT...",
	Short: "Remove one or more contexts (and their cluster + user) from a kubeconfig",
	Example: `  # Remove a single context from the active kubeconfig
  ktidy remove old-cluster

  # Remove multiple contexts at once
  ktidy remove staging dev old-cluster

  # Preview what would be removed without writing anything
  ktidy remove --dry-run prod`,
	Args:  cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		path, _ := cmd.Flags().GetString("kubeconfig")
		path = resolveKubeconfig(path)
		cfg, err := kubeconfig.Load(path)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		names := make([]string, 0, len(cfg.Contexts))
		for name := range cfg.Contexts {
			names = append(names, name)
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: runRemove,
}

var removeKubeconfig string

func init() {
	removeCmd.Flags().StringVarP(&removeKubeconfig, "kubeconfig", "k", "", "path to kubeconfig (default: $KUBECONFIG or ~/.kube/config)")
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
	path := resolveKubeconfig(removeKubeconfig)
	cfg, err := kubeconfig.Load(path)
	if err != nil {
		return fmt.Errorf("loading %s: %w", path, err)
	}

	missing := kubeconfig.Remove(cfg, args)
	for _, name := range missing {
		warnf(cmd, "context %q not found", name)
	}

	if dryRunSkip(cmd, "would write updated config to %s", path) {
		return nil
	}
	return kubeconfig.Write(cfg, path)
}

