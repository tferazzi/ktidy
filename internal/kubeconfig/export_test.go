package kubeconfig

import (
	"testing"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func makeFullConfig(ctxName, clusterName, userName string) *clientcmdapi.Config {
	cfg := clientcmdapi.NewConfig()
	cfg.Contexts[ctxName] = &clientcmdapi.Context{Cluster: clusterName, AuthInfo: userName}
	cfg.Clusters[clusterName] = &clientcmdapi.Cluster{Server: "https://example.com"}
	cfg.AuthInfos[userName] = &clientcmdapi.AuthInfo{Token: "tok"}
	return cfg
}

func TestExport(t *testing.T) {
	src := makeFullConfig("ctx-a", "cluster-a", "user-a")
	// add a second context that should NOT appear in output
	src.Contexts["ctx-b"] = &clientcmdapi.Context{Cluster: "cluster-b", AuthInfo: "user-b"}
	src.Clusters["cluster-b"] = &clientcmdapi.Cluster{}
	src.AuthInfos["user-b"] = &clientcmdapi.AuthInfo{}

	tests := []struct {
		name         string
		contexts     []string
		wantErr      bool
		wantContexts int
		wantCurrent  string
	}{
		{"single context", []string{"ctx-a"}, false, 1, "ctx-a"},
		{"multiple contexts", []string{"ctx-a", "ctx-b"}, false, 2, ""},
		{"missing context", []string{"ctx-missing"}, true, 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := Export(src, tt.contexts)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if len(out.Contexts) != tt.wantContexts {
				t.Errorf("contexts: got %d, want %d", len(out.Contexts), tt.wantContexts)
			}
			if out.CurrentContext != tt.wantCurrent {
				t.Errorf("currentContext: got %q, want %q", out.CurrentContext, tt.wantCurrent)
			}
			// cluster and user referenced by ctx-a must be present
			if _, ok := out.Clusters["cluster-a"]; !ok {
				t.Error("cluster-a missing from export")
			}
			if _, ok := out.AuthInfos["user-a"]; !ok {
				t.Error("user-a missing from export")
			}
		})
	}
}
