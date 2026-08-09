package kubeconfig

import (
	"testing"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func makeConfig(contexts, clusters, users []string) *clientcmdapi.Config {
	cfg := clientcmdapi.NewConfig()
	for _, name := range contexts {
		cfg.Contexts[name] = &clientcmdapi.Context{}
	}
	for _, name := range clusters {
		cfg.Clusters[name] = &clientcmdapi.Cluster{}
	}
	for _, name := range users {
		cfg.AuthInfos[name] = &clientcmdapi.AuthInfo{}
	}
	return cfg
}

func TestDetect(t *testing.T) {
	a := makeConfig([]string{"ctx-a", "shared"}, []string{"cluster-a"}, []string{"user-a"})
	b := makeConfig([]string{"ctx-b", "shared"}, []string{"cluster-a"}, []string{"user-b"})

	c := Detect([]*clientcmdapi.Config{a, b})

	if !c.Any() {
		t.Fatal("expected conflicts, got none")
	}
	if len(c.Contexts) != 1 || c.Contexts[0] != "shared" {
		t.Errorf("contexts: got %v, want [shared]", c.Contexts)
	}
	if len(c.Clusters) != 1 || c.Clusters[0] != "cluster-a" {
		t.Errorf("clusters: got %v, want [cluster-a]", c.Clusters)
	}
	if len(c.Users) != 0 {
		t.Errorf("users: got %v, want []", c.Users)
	}
}

func TestMerge_lastWriteWins(t *testing.T) {
	a := makeConfig([]string{"ctx"}, nil, nil)
	a.Contexts["ctx"].Namespace = "ns-a"

	b := makeConfig([]string{"ctx"}, nil, nil)
	b.Contexts["ctx"].Namespace = "ns-b"

	merged, err := Merge([]*clientcmdapi.Config{a, b}, false)
	if err != nil {
		t.Fatal(err)
	}
	if merged.Contexts["ctx"].Namespace != "ns-b" {
		t.Errorf("got %q, want ns-b", merged.Contexts["ctx"].Namespace)
	}
}

func TestDetectAgainst(t *testing.T) {
	base := makeConfig([]string{"ctx-a"}, []string{"cluster-a"}, []string{"user-a"})
	incoming := makeConfig([]string{"ctx-a", "ctx-b"}, []string{"cluster-b"}, []string{"user-a"})

	c := DetectAgainst(base, []*clientcmdapi.Config{incoming})

	if len(c.Contexts) != 1 || c.Contexts[0] != "ctx-a" {
		t.Errorf("contexts: got %v, want [ctx-a]", c.Contexts)
	}
	if len(c.Clusters) != 0 {
		t.Errorf("clusters: got %v, want []", c.Clusters)
	}
	if len(c.Users) != 1 || c.Users[0] != "user-a" {
		t.Errorf("users: got %v, want [user-a]", c.Users)
	}
}

func TestMerge_noConflict(t *testing.T) {
	a := makeConfig([]string{"ctx-a"}, []string{"cluster-a"}, []string{"user-a"})
	b := makeConfig([]string{"ctx-b"}, []string{"cluster-b"}, []string{"user-b"})

	merged, err := Merge([]*clientcmdapi.Config{a, b}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(merged.Contexts) != 2 {
		t.Errorf("got %d contexts, want 2", len(merged.Contexts))
	}
}
