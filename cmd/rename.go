package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tferazzi/ktidy/internal/kubeconfig"
)

var renameCmd = &cobra.Command{
	Use:   "rename <OLD> <NEW>",
	Short: "Rename a context in a kubeconfig file",
	Args:  cobra.ExactArgs(2),
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
		fmt.Fprintf(cmd.ErrOrStderr(), "warning: active context renamed to %q\n", args[1])
	}

	return kubeconfig.Write(cfg, path)
}
