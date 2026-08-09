package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tferazzi/ktidy/internal/kubeconfig"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var splitCmd = &cobra.Command{
	Use:   "split [flags] CONTEXT",
	Short: "Export a context to stdout; with --remove also delete it from the source file",
	Example: `  # Export a context (identical to ktidy export for a single context)
  ktidy split prod > prod.yaml

  # Export and remove the context from the source file
  ktidy split --remove prod > prod.yaml

  # Preview what --remove would delete without touching the file
  ktidy split --remove --dry-run prod`,
	Args:  cobra.ExactArgs(1),
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
	RunE: runSplit,
}

var splitRemove bool

func init() {
	splitCmd.Flags().StringP("kubeconfig", "k", "", "path to kubeconfig file (default: $KUBECONFIG or ~/.kube/config)")
	splitCmd.Flags().BoolVar(&splitRemove, "remove", false, "remove the context from the source file after exporting")
	rootCmd.AddCommand(splitCmd)
}

func runSplit(cmd *cobra.Command, args []string) error {
	contextName := args[0]
	path, _ := cmd.Flags().GetString("kubeconfig")
	resolvedPath := resolveKubeconfig(path)

	src, err := kubeconfig.Load(resolvedPath)
	if err != nil {
		return err
	}

	out, err := buildSplitOutput(src, contextName)
	if err != nil {
		return err
	}

	b, err := clientcmd.Write(*out)
	if err != nil {
		return err
	}
	if _, err := os.Stdout.Write(b); err != nil {
		return err
	}

	if splitRemove {
		missing := kubeconfig.Remove(src, []string{contextName})
		if len(missing) > 0 {
			return fmt.Errorf("context %q not found", contextName)
		}
		if dryRunSkip(cmd, "would remove context %q from %s", contextName, resolvedPath) {
			return nil
		}
		return kubeconfig.Write(src, resolvedPath)
	}
	return nil
}

func buildSplitOutput(src *clientcmdapi.Config, contextName string) (*clientcmdapi.Config, error) {
	ctx, ok := src.Contexts[contextName]
	if !ok {
		return nil, fmt.Errorf("context %q not found", contextName)
	}
	out := clientcmdapi.NewConfig()
	out.Contexts[contextName] = ctx
	out.CurrentContext = contextName
	if cluster, ok := src.Clusters[ctx.Cluster]; ok {
		out.Clusters[ctx.Cluster] = cluster
	}
	if user, ok := src.AuthInfos[ctx.AuthInfo]; ok {
		out.AuthInfos[ctx.AuthInfo] = user
	}
	if err := clientcmdapi.FlattenConfig(out); err != nil {
		return nil, err
	}
	return out, nil
}
