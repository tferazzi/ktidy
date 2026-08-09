package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/tferazzi/ktidy/internal/kubeconfig"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var importCmd = &cobra.Command{
	Use:   "import [flags] CONFIG...",
	Short: "Import kubeconfig(s) into the current kubeconfig",
	Example: `  # Preview what would be imported without writing anything
  ktidy import --dry-run new-cluster.yaml

  # Import and save back to ~/.kube/config
  ktidy import --save new-cluster.yaml

  # Overwrite existing conflicting entries
  ktidy import --save --force updated-cluster.yaml

  # Import from stdin (e.g. from a cloud provider command)
  aws eks update-kubeconfig --name my-cluster --dry-run | ktidy import --stdin --save`,
	RunE: runImport,
}

var (
	importPreserveStructure bool
	importSave              bool
	importStdin             bool
	importForce             bool
)

func init() {
	importCmd.Flags().BoolVarP(&importPreserveStructure, "preserve-structure", "p", false, "keep credential files as external paths instead of inlining")
	importCmd.Flags().BoolVarP(&importSave, "save", "s", false, "write result back to the resolved config path")
	importCmd.Flags().BoolVarP(&importStdin, "stdin", "i", false, "read config from stdin")
	importCmd.Flags().BoolVarP(&importForce, "force", "f", false, "overwrite conflicting entries instead of skipping them")
	importCmd.MarkFlagsMutuallyExclusive("stdin", "save")
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) error {
	if !importStdin && len(args) == 0 {
		return fmt.Errorf("at least one CONFIG file is required, or use --stdin")
	}

	incoming := make([]*clientcmdapi.Config, 0, len(args)+1)

	if importStdin {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("reading stdin: %w", err)
		}
		cfg, err := clientcmd.Load(data)
		if err != nil {
			return fmt.Errorf("parsing stdin: %w", err)
		}
		incoming = append(incoming, cfg)
	}

	for _, path := range args {
		cfg, err := kubeconfig.Load(path)
		if err != nil {
			return fmt.Errorf("loading %s: %w", path, err)
		}
		incoming = append(incoming, cfg)
	}

	destPath := kubeconfig.ResolveConfigPath()

	var base *clientcmdapi.Config
	if _, err := os.Stat(destPath); err == nil {
		base, err = kubeconfig.Load(destPath)
		if err != nil {
			return fmt.Errorf("loading %s: %w", destPath, err)
		}
	} else {
		base = clientcmdapi.NewConfig()
	}

	conflicts := kubeconfig.DetectAgainst(base, incoming)
	if conflicts.Any() && !importForce {
		for _, name := range conflicts.Contexts {
			warnf(cmd, "skipping context %q: already exists (use --force to overwrite)", name)
		}
		for _, name := range conflicts.Clusters {
			warnf(cmd, "skipping cluster %q: already exists (use --force to overwrite)", name)
		}
		for _, name := range conflicts.Users {
			warnf(cmd, "skipping user %q: already exists (use --force to overwrite)", name)
		}
		// filter out conflicting entries from incoming before merging
		incoming = filterConflicts(incoming, conflicts)
	}

	all := append([]*clientcmdapi.Config{base}, incoming...)
	merged, err := kubeconfig.Merge(all, importPreserveStructure)
	if err != nil {
		return err
	}

	out, err := clientcmd.Write(*merged)
	if err != nil {
		return err
	}

	if importSave {
		if dryRunSkip(cmd, "would write merged config to %s", destPath) {
			return nil
		}
		return kubeconfig.Write(merged, destPath)
	}
	_, err = os.Stdout.Write(out)
	return err
}

func filterConflicts(configs []*clientcmdapi.Config, c kubeconfig.Conflicts) []*clientcmdapi.Config {
	ctxSet := toSet(c.Contexts)
	clsSet := toSet(c.Clusters)
	usrSet := toSet(c.Users)

	result := make([]*clientcmdapi.Config, 0, len(configs))
	for _, cfg := range configs {
		clean := clientcmdapi.NewConfig()
		for name, v := range cfg.Contexts {
			if !ctxSet[name] {
				clean.Contexts[name] = v
			}
		}
		for name, v := range cfg.Clusters {
			if !clsSet[name] {
				clean.Clusters[name] = v
			}
		}
		for name, v := range cfg.AuthInfos {
			if !usrSet[name] {
				clean.AuthInfos[name] = v
			}
		}
		if cfg.CurrentContext != "" && !ctxSet[cfg.CurrentContext] {
			clean.CurrentContext = cfg.CurrentContext
		}
		result = append(result, clean)
	}
	return result
}

func toSet(names []string) map[string]bool {
	s := make(map[string]bool, len(names))
	for _, n := range names {
		s[n] = true
	}
	return s
}
