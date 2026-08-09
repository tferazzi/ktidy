package kubeconfig

import (
	"maps"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Merge combines multiple configs into one. Later configs win on conflict.
// Pass preserveStructure=true to keep external credential file references
// instead of inlining their contents.
func Merge(configs []*clientcmdapi.Config, preserveStructure bool) (*clientcmdapi.Config, error) {
	out := clientcmdapi.NewConfig()

	for _, cfg := range configs {
		mergeInto(out, cfg)
	}

	if preserveStructure {
		return out, clientcmdapi.MinifyConfig(out)
	}
	return out, clientcmdapi.FlattenConfig(out)
}

func mergeInto(dst, src *clientcmdapi.Config) {
	maps.Copy(dst.Contexts, src.Contexts)
	maps.Copy(dst.Clusters, src.Clusters)
	maps.Copy(dst.AuthInfos, src.AuthInfos)
	if src.CurrentContext != "" && dst.CurrentContext == "" {
		dst.CurrentContext = src.CurrentContext
	}
}
