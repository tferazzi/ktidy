package kubeconfig

import (
	"testing"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func TestRemove(t *testing.T) {
	cfg := func() *clientcmdapi.Config {
		c := clientcmdapi.NewConfig()
		c.Contexts["dev"] = &clientcmdapi.Context{Cluster: "dev-cluster", AuthInfo: "dev-user"}
		c.Clusters["dev-cluster"] = &clientcmdapi.Cluster{}
		c.AuthInfos["dev-user"] = &clientcmdapi.AuthInfo{}
		c.CurrentContext = "dev"
		return c
	}

	t.Run("removes context cluster and user", func(t *testing.T) {
		c := cfg()
		if err := Remove(c, "dev"); err != nil {
			t.Fatal(err)
		}
		if _, ok := c.Contexts["dev"]; ok {
			t.Error("context still present")
		}
		if _, ok := c.Clusters["dev-cluster"]; ok {
			t.Error("cluster still present")
		}
		if _, ok := c.AuthInfos["dev-user"]; ok {
			t.Error("user still present")
		}
		if c.CurrentContext != "" {
			t.Errorf("current-context not cleared, got %q", c.CurrentContext)
		}
	})

	t.Run("error on missing context", func(t *testing.T) {
		if err := Remove(cfg(), "nope"); err == nil {
			t.Error("expected error, got nil")
		}
	})
}
