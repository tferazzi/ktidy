package kubeconfig

import (
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func Load(path string) (*clientcmdapi.Config, error) {
	return clientcmd.LoadFromFile(path)
}

func Write(cfg *clientcmdapi.Config, path string) error {
	return clientcmd.WriteToFile(*cfg, path)
}
