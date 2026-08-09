package kubeconfig

import (
	"fmt"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Split extracts the named context and its referenced cluster + user into a
// minimal config. Returns an error if the context does not exist.
func Split(cfg *clientcmdapi.Config, contextName string) (*clientcmdapi.Config, error) {
	ctx, ok := cfg.Contexts[contextName]
	if !ok {
		return nil, fmt.Errorf("context %q not found", contextName)
	}

	out := clientcmdapi.NewConfig()
	out.Contexts[contextName] = ctx
	out.CurrentContext = contextName
	if cluster, ok := cfg.Clusters[ctx.Cluster]; ok {
		out.Clusters[ctx.Cluster] = cluster
	}
	if user, ok := cfg.AuthInfos[ctx.AuthInfo]; ok {
		out.AuthInfos[ctx.AuthInfo] = user
	}
	return out, nil
}
