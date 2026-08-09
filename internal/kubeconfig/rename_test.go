package kubeconfig

import (
	"testing"
)

func TestRename(t *testing.T) {
	tests := []struct {
		name       string
		current    string
		old, new   string
		wantWarn   bool
		wantErr    bool
	}{
		{"basic rename", "", "ctx-a", "ctx-b", false, false},
		{"renames current context", "ctx-a", "ctx-a", "ctx-b", true, false},
		{"old not found", "", "missing", "ctx-b", false, true},
		{"new already exists", "", "ctx-a", "ctx-c", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := makeConfig([]string{"ctx-a", "ctx-c"}, nil, nil)
			cfg.CurrentContext = tt.current

			warned, err := Rename(cfg, tt.old, tt.new)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err=%v, wantErr=%v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if warned != tt.wantWarn {
				t.Errorf("warned=%v, want %v", warned, tt.wantWarn)
			}
			if _, ok := cfg.Contexts[tt.new]; !ok {
				t.Errorf("new context %q not found", tt.new)
			}
			if _, ok := cfg.Contexts[tt.old]; ok {
				t.Errorf("old context %q still present", tt.old)
			}
			if tt.wantWarn && cfg.CurrentContext != tt.new {
				t.Errorf("currentContext=%q, want %q", cfg.CurrentContext, tt.new)
			}
		})
	}
}
