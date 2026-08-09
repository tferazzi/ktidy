package kubeconfig

import (
	"fmt"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Export extracts the named contexts (plus their referenced cluster and user)
// from src into a minimal kubeconfig.
func Export(src *clientcmdapi.Config, contextNames []string) (*clientcmdapi.Config, error) {
	out := clientcmdapi.NewConfig()
	for _, name := range contextNames {
		ctx, ok := src.Contexts[name]
		if !ok {
			return nil, fmt.Errorf("context not found: %s", name)
		}
		out.Contexts[name] = ctx
		if ctx.Cluster != "" {
			if cl, ok := src.Clusters[ctx.Cluster]; ok {
				out.Clusters[ctx.Cluster] = cl
			}
		}
		if ctx.AuthInfo != "" {
			if u, ok := src.AuthInfos[ctx.AuthInfo]; ok {
				out.AuthInfos[ctx.AuthInfo] = u
			}
		}
	}
	if len(contextNames) == 1 {
		out.CurrentContext = contextNames[0]
	}
	return out, nil
}
