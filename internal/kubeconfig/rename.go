package kubeconfig

import (
	"fmt"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Rename re-keys a context entry from oldName to newName.
// Returns (warnRenamed, error) where warnRenamed is true when currentContext was updated.
func Rename(cfg *clientcmdapi.Config, oldName, newName string) (bool, error) {
	if _, ok := cfg.Contexts[oldName]; !ok {
		return false, fmt.Errorf("context %q not found", oldName)
	}
	if _, ok := cfg.Contexts[newName]; ok {
		return false, fmt.Errorf("context %q already exists", newName)
	}

	cfg.Contexts[newName] = cfg.Contexts[oldName]
	delete(cfg.Contexts, oldName)

	if cfg.CurrentContext == oldName {
		cfg.CurrentContext = newName
		return true, nil
	}
	return false, nil
}
