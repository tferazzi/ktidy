package kubeconfig

import (
	"testing"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func makeConfigFromMap(contexts map[string][2]string) *clientcmdapi.Config {
	cfg := clientcmdapi.NewConfig()
	for ctxName, refs := range contexts {
		cfg.Contexts[ctxName] = &clientcmdapi.Context{Cluster: refs[0], AuthInfo: refs[1]}
		cfg.Clusters[refs[0]] = &clientcmdapi.Cluster{}
		cfg.AuthInfos[refs[1]] = &clientcmdapi.AuthInfo{}
	}
	return cfg
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name         string
		setup        map[string][2]string
		current      string
		remove       []string
		wantMissing  []string
		wantContexts []string
		wantCurrent  string
	}{
		{
			name:         "removes context cluster and user",
			setup:        map[string][2]string{"ctx-a": {"cluster-a", "user-a"}, "ctx-b": {"cluster-b", "user-b"}},
			remove:       []string{"ctx-a"},
			wantMissing:  nil,
			wantContexts: []string{"ctx-b"},
		},
		{
			name:        "warns on missing context",
			setup:       map[string][2]string{"ctx-a": {"cluster-a", "user-a"}},
			remove:      []string{"ctx-missing"},
			wantMissing: []string{"ctx-missing"},
			wantContexts: []string{"ctx-a"},
		},
		{
			name:         "clears current context when removed",
			setup:        map[string][2]string{"ctx-a": {"cluster-a", "user-a"}},
			current:      "ctx-a",
			remove:       []string{"ctx-a"},
			wantMissing:  nil,
			wantContexts: []string{},
			wantCurrent:  "",
		},
		{
			name:         "removes multiple contexts",
			setup:        map[string][2]string{"ctx-a": {"cluster-a", "user-a"}, "ctx-b": {"cluster-b", "user-b"}, "ctx-c": {"cluster-c", "user-c"}},
			remove:       []string{"ctx-a", "ctx-b"},
			wantMissing:  nil,
			wantContexts: []string{"ctx-c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := makeConfigFromMap(tt.setup)
			cfg.CurrentContext = tt.current

			missing := Remove(cfg, tt.remove)

			if len(missing) != len(tt.wantMissing) {
				t.Errorf("missing: got %v, want %v", missing, tt.wantMissing)
			}
			if len(cfg.Contexts) != len(tt.wantContexts) {
				t.Errorf("contexts count: got %d, want %d", len(cfg.Contexts), len(tt.wantContexts))
			}
			if cfg.CurrentContext != tt.wantCurrent {
				t.Errorf("current context: got %q, want %q", cfg.CurrentContext, tt.wantCurrent)
			}
			// verify removed cluster/user entries are gone
			for _, name := range tt.remove {
				if ctx, ok := tt.setup[name]; ok {
					if _, exists := cfg.Clusters[ctx[0]]; exists {
						t.Errorf("cluster %q should be deleted", ctx[0])
					}
					if _, exists := cfg.AuthInfos[ctx[1]]; exists {
						t.Errorf("user %q should be deleted", ctx[1])
					}
				}
			}
		})
	}
}
