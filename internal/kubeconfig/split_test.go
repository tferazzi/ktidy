package kubeconfig

import (
	"testing"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func TestSplit(t *testing.T) {
	cfg := clientcmdapi.NewConfig()
	cfg.Contexts["dev"] = &clientcmdapi.Context{Cluster: "dev-cluster", AuthInfo: "dev-user"}
	cfg.Contexts["prod"] = &clientcmdapi.Context{Cluster: "prod-cluster", AuthInfo: "prod-user"}
	cfg.Clusters["dev-cluster"] = &clientcmdapi.Cluster{Server: "https://dev"}
	cfg.Clusters["prod-cluster"] = &clientcmdapi.Cluster{Server: "https://prod"}
	cfg.AuthInfos["dev-user"] = &clientcmdapi.AuthInfo{Token: "dev-tok"}
	cfg.AuthInfos["prod-user"] = &clientcmdapi.AuthInfo{Token: "prod-tok"}

	t.Run("extracts context cluster and user", func(t *testing.T) {
		out, err := Split(cfg, "dev")
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := out.Contexts["dev"]; !ok {
			t.Error("dev context missing")
		}
		if _, ok := out.Clusters["dev-cluster"]; !ok {
			t.Error("dev-cluster missing")
		}
		if _, ok := out.AuthInfos["dev-user"]; !ok {
			t.Error("dev-user missing")
		}
		if _, ok := out.Contexts["prod"]; ok {
			t.Error("prod context should not be present")
		}
		if out.CurrentContext != "dev" {
			t.Errorf("currentContext: got %q, want dev", out.CurrentContext)
		}
	})

	t.Run("error on missing context", func(t *testing.T) {
		if _, err := Split(cfg, "nope"); err == nil {
			t.Error("expected error, got nil")
		}
	})
}
