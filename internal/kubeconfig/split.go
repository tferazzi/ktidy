package kubeconfig

import (
	"fmt"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Remove deletes a context and its associated cluster and user from cfg.
// Returns an error if the context does not exist.
func Remove(cfg *clientcmdapi.Config, contextName string) error {
	ctx, ok := cfg.Contexts[contextName]
	if !ok {
		return fmt.Errorf("context %q not found", contextName)
	}
	delete(cfg.Clusters, ctx.Cluster)
	delete(cfg.AuthInfos, ctx.AuthInfo)
	delete(cfg.Contexts, contextName)
	if cfg.CurrentContext == contextName {
		cfg.CurrentContext = ""
	}
	return nil
}
