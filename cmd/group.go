package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tferazzi/ktidy/internal/kubeconfig"
)

var groupsDir string

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage kubeconfig groups",
	Example: `  # One-time shell setup — add to ~/.zshrc or ~/.bashrc
  eval "$(ktidy group activate)"`,
}

var groupListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List groups and their contexts",
	Example: `  ktidy group list`,
	Args:    cobra.NoArgs,
	RunE:    runGroupList,
}

var groupCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new group directory",
	Example: `  ktidy group create aws
  ktidy group create gcp`,
	Args:              cobra.ExactArgs(1),
	RunE:              runGroupCreate,
	ValidArgsFunction: cobra.NoFileCompletions,
}

var groupAddMove bool

var groupAddCmd = &cobra.Command{
	Use:   "add <group> <file>",
	Short: "Add a kubeconfig file to a group",
	Example: `  # Copy a kubeconfig into the aws group
  ktidy group add aws ~/Downloads/eks-prod.yaml

  # Move instead of copy
  ktidy group add --move aws ~/Downloads/eks-staging.yaml`,
	Args: cobra.ExactArgs(2),
	RunE: runGroupAdd,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return groupNames(), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveDefault
	},
}

var groupActivateCmd = &cobra.Command{
	Use:   "activate [<group>...]",
	Short: "Print export KUBECONFIG=... for eval",
	Example: `  # Activate all groups (add to ~/.zshrc)
  eval "$(ktidy group activate)"

  # Activate only the aws and gcp groups
  eval "$(ktidy group activate aws gcp)"`,
	RunE: runGroupActivate,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return groupNames(), cobra.ShellCompDirectiveNoFileComp
	},
}

var groupRemoveForce bool

var groupRemoveCmd = &cobra.Command{
	Use:   "remove <group>",
	Short: "Remove a group directory",
	Example: `  # Remove with confirmation prompt
  ktidy group remove old-group

  # Remove without prompting
  ktidy group remove --force old-group`,
	Args: cobra.ExactArgs(1),
	RunE: runGroupRemove,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return groupNames(), cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	defaultGroupsDir := os.Getenv("KTIDY_GROUPS_DIR")
	if defaultGroupsDir == "" {
		home, _ := os.UserHomeDir()
		defaultGroupsDir = filepath.Join(home, ".kube", "groups")
	}

	groupCmd.PersistentFlags().StringVar(&groupsDir, "groups-dir", defaultGroupsDir, "root directory for kubeconfig groups ($KTIDY_GROUPS_DIR)")
	groupAddCmd.Flags().BoolVarP(&groupAddMove, "move", "m", false, "move the file instead of copying it")
	groupRemoveCmd.Flags().BoolVarP(&groupRemoveForce, "force", "f", false, "skip confirmation prompt")

	groupCmd.AddCommand(groupListCmd, groupCreateCmd, groupAddCmd, groupActivateCmd, groupRemoveCmd)
	rootCmd.AddCommand(groupCmd)
}

func runGroupList(cmd *cobra.Command, _ []string) error {
	entries, err := os.ReadDir(groupsDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", e.Name())
		files, err := kubeconfigFiles(filepath.Join(groupsDir, e.Name()))
		if err != nil {
			return err
		}
		for _, f := range files {
			cfg, err := kubeconfig.Load(f)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "  warning: cannot read %s: %v\n", f, err)
				continue
			}
			for name := range cfg.Contexts {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", name)
			}
		}
	}
	return nil
}

func runGroupCreate(cmd *cobra.Command, args []string) error {
	dir := filepath.Join(groupsDir, args[0])
	if dryRunSkip(cmd, "would create group directory %s", dir) {
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "created group %q\n", args[0])
	return nil
}

func runGroupAdd(cmd *cobra.Command, args []string) error {
	group, src := args[0], args[1]
	groupDir := filepath.Join(groupsDir, group)

	if err := os.MkdirAll(groupDir, 0o755); err != nil {
		return err
	}

	incomingCfg, err := kubeconfig.Load(src)
	if err != nil {
		return fmt.Errorf("loading %s: %w", src, err)
	}

	// warn on context name conflicts with existing files in the group
	existing, err := kubeconfigFiles(groupDir)
	if err != nil {
		return err
	}
	for _, f := range existing {
		existingCfg, err := kubeconfig.Load(f)
		if err != nil {
			continue
		}
		for name := range incomingCfg.Contexts {
			if _, ok := existingCfg.Contexts[name]; ok {
				fmt.Fprintf(cmd.ErrOrStderr(), "warning: context %q already exists in %s\n", name, f)
			}
		}
	}

	dst := filepath.Join(groupDir, filepath.Base(src))
	if dryRunSkip(cmd, "would %s %s to %s", map[bool]string{true: "move", false: "copy"}[groupAddMove], src, dst) {
		return nil
	}
	if groupAddMove {
		if err := os.Rename(src, dst); err != nil {
			return err
		}
	} else {
		if err := copyFile(src, dst); err != nil {
			return err
		}
	}
	fmt.Fprintf(cmd.OutOrStdout(), "added %s to group %q\n", filepath.Base(src), group)
	return nil
}

func runGroupActivate(cmd *cobra.Command, args []string) error {
	var dirs []string
	if len(args) == 0 {
		entries, err := os.ReadDir(groupsDir)
		if os.IsNotExist(err) {
			fmt.Fprintf(cmd.OutOrStdout(), "export KUBECONFIG=\n")
			return nil
		}
		if err != nil {
			return err
		}
		for _, e := range entries {
			if e.IsDir() {
				dirs = append(dirs, filepath.Join(groupsDir, e.Name()))
			}
		}
	} else {
		for _, name := range args {
			dirs = append(dirs, filepath.Join(groupsDir, name))
		}
	}

	var files []string
	for _, d := range dirs {
		fs, err := kubeconfigFiles(d)
		if err != nil {
			return fmt.Errorf("reading group %s: %w", filepath.Base(d), err)
		}
		files = append(files, fs...)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "export KUBECONFIG=%s\n", strings.Join(files, ":"))
	return nil
}

func runGroupRemove(cmd *cobra.Command, args []string) error {
	group := args[0]
	dir := filepath.Join(groupsDir, group)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("group %q does not exist", group)
	}

	if !groupRemoveForce {
		fmt.Fprintf(cmd.OutOrStdout(), "remove group %q and all its files? [y/N] ", group)
		scanner := bufio.NewScanner(cmd.InOrStdin())
		scanner.Scan()
		if strings.ToLower(strings.TrimSpace(scanner.Text())) != "y" {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
	}

	if dryRunSkip(cmd, "would remove group directory %s", dir) {
		return nil
	}
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "removed group %q\n", group)
	return nil
}

// kubeconfigFiles returns all .yaml/.yml files in dir.
func kubeconfigFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
			out = append(out, filepath.Join(dir, name))
		}
	}
	return out, nil
}

// groupNames returns the list of existing group names for shell completions.
func groupNames() []string {
	entries, err := os.ReadDir(groupsDir)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}
