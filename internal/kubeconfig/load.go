package kubeconfig

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func Load(path string) (*clientcmdapi.Config, error) {
	return clientcmd.LoadFromFile(path)
}

func Write(cfg *clientcmdapi.Config, path string) error {
	return clientcmd.WriteToFile(*cfg, path)
}

// ResolveConfigPath returns the kubeconfig path to use:
// KUBECONFIG env → XDG_CONFIG_HOME/kube/config → ~/.kube/config
func ResolveConfigPath() string {
	if p := os.Getenv("KUBECONFIG"); p != "" {
		return p
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "kube", "config")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".kube", "config")
}
