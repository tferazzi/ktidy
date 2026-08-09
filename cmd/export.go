package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tferazzi/ktidy/internal/kubeconfig"
	"k8s.io/client-go/tools/clientcmd"
)

var exportCmd = &cobra.Command{
	Use:   "export [flags] CONTEXT...",
	Short: "Extract one or more contexts into a minimal kubeconfig and print to stdout",
	Example: `  # Export a single context to a file
  ktidy export prod > prod.yaml

  # Export two contexts from a specific file
  ktidy export -k ~/other.yaml staging prod > two-contexts.yaml`,
	Args: cobra.MinimumNArgs(1),
	RunE: runExport,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		path := resolveKubeconfig(exportKubeconfig)
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
}

var exportKubeconfig string

func init() {
	exportCmd.Flags().StringVarP(&exportKubeconfig, "kubeconfig", "k", "", "kubeconfig file (default: $KUBECONFIG or ~/.kube/config); comma-delimited or repeated")
	rootCmd.AddCommand(exportCmd)
}

func resolveKubeconfig(flag string) string {
	if flag != "" {
		// take first path if comma-delimited
		return strings.SplitN(flag, ",", 2)[0]
	}
	if env := os.Getenv("KUBECONFIG"); env != "" {
		return strings.SplitN(env, ",", 2)[0]
	}
	return clientcmd.RecommendedHomeFile
}

func runExport(cmd *cobra.Command, args []string) error {
	path := resolveKubeconfig(exportKubeconfig)
	cfg, err := kubeconfig.Load(path)
	if err != nil {
		return fmt.Errorf("loading %s: %w", path, err)
	}

	exported, err := kubeconfig.Export(cfg, args)
	if err != nil {
		return err
	}

	out, err := clientcmd.Write(*exported)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(out)
	return err
}
