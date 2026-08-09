package kubeconfig

import (
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Remove deletes the named contexts and their referenced clusters and users
// from cfg. Returns the names of any contexts that were not found.
func Remove(cfg *clientcmdapi.Config, contexts []string) []string {
	var missing []string
	for _, name := range contexts {
		ctx, ok := cfg.Contexts[name]
		if !ok {
			missing = append(missing, name)
			continue
		}
		delete(cfg.Clusters, ctx.Cluster)
		delete(cfg.AuthInfos, ctx.AuthInfo)
		delete(cfg.Contexts, name)
		if cfg.CurrentContext == name {
			cfg.CurrentContext = ""
		}
	}
	return missing
}
