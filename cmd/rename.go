package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tferazzi/ktidy/internal/kubeconfig"
)

var renameCmd = &cobra.Command{
	Use:   "rename <OLD> <NEW>",
	Short: "Rename a context in a kubeconfig file",
	Example: `  # Rename a verbose ARN context to something human-friendly
  ktidy rename arn:aws:eks:eu-west-1:123456789:cluster/prod prod

  # Rename in a specific file
  ktidy rename -k ~/other.yaml old-name new-name

  # Preview the rename without writing anything
  ktidy rename --dry-run old-name new-name`,
	Args: cobra.ExactArgs(2),
	RunE:  runRename,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		path, _ := cmd.Flags().GetString("kubeconfig")
		cfg, err := kubeconfig.Load(resolveKubeconfig(path))
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		names := make([]string, 0, len(cfg.Contexts))
		for name := range cfg.Contexts {
			names = append(names, name)
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	},
}

var renameKubeconfig string

func init() {
	renameCmd.Flags().StringVarP(&renameKubeconfig, "kubeconfig", "k", "", "path to kubeconfig file")
	rootCmd.AddCommand(renameCmd)
}

func runRename(cmd *cobra.Command, args []string) error {
	path := resolveKubeconfig(renameKubeconfig)
	cfg, err := kubeconfig.Load(path)
	if err != nil {
		return fmt.Errorf("loading %s: %w", path, err)
	}

	warned, err := kubeconfig.Rename(cfg, args[0], args[1])
	if err != nil {
		return err
	}
	if warned {
		warnf(cmd, "active context renamed to %q", args[1])
	}

	if dryRunSkip(cmd, "would rename context %q to %q in %s", args[0], args[1], path) {
		return nil
	}
	return kubeconfig.Write(cfg, path)
}
