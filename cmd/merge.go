package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tferazzi/ktidy/internal/kubeconfig"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var mergeCmd = &cobra.Command{
	Use:   "merge [flags] CONFIG...",
	Short: "Merge multiple kubeconfig files and print to stdout",
	Example: `  # Merge two files and write the result to ~/.kube/config
  ktidy merge work.yaml personal.yaml > ~/.kube/config

  # Merge while keeping credential files as external references
  ktidy merge -p cluster-a.yaml cluster-b.yaml > merged.yaml`,
	Args: cobra.MinimumNArgs(1),
	RunE: runMerge,
}

var mergePreserveStructure bool

func init() {
	mergeCmd.Flags().BoolVarP(&mergePreserveStructure, "preserve-structure", "p", false, "keep credential files as external paths instead of inlining")
	rootCmd.AddCommand(mergeCmd)
}

func runMerge(cmd *cobra.Command, args []string) error {
	configs := make([]*clientcmdapi.Config, 0, len(args))
	for _, path := range args {
		cfg, err := kubeconfig.Load(path)
		if err != nil {
			return fmt.Errorf("loading %s: %w", path, err)
		}
		configs = append(configs, cfg)
	}

	if conflicts := kubeconfig.Detect(configs); conflicts.Any() {
		printConflicts(cmd, conflicts)
	}

	merged, err := kubeconfig.Merge(configs, mergePreserveStructure)
	if err != nil {
		return err
	}

	out, err := clientcmd.Write(*merged)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(out)
	return err
}

func printConflicts(cmd *cobra.Command, c kubeconfig.Conflicts) {
	for _, name := range c.Contexts {
		warnf(cmd, "duplicate context %q — last-write wins", name)
	}
	for _, name := range c.Clusters {
		warnf(cmd, "duplicate cluster %q — last-write wins", name)
	}
	for _, name := range c.Users {
		warnf(cmd, "duplicate user %q — last-write wins", name)
	}
}
