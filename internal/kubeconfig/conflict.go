package kubeconfig

import clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

// Conflicts holds the names that appear in more than one config.
type Conflicts struct {
	Contexts []string
	Clusters []string
	Users    []string
}

func (c Conflicts) Any() bool {
	return len(c.Contexts) > 0 || len(c.Clusters) > 0 || len(c.Users) > 0
}

// Detect returns the names that collide across configs.
func Detect(configs []*clientcmdapi.Config) Conflicts {
	contexts := map[string]int{}
	clusters := map[string]int{}
	users := map[string]int{}

	for _, cfg := range configs {
		for name := range cfg.Contexts {
			contexts[name]++
		}
		for name := range cfg.Clusters {
			clusters[name]++
		}
		for name := range cfg.AuthInfos {
			users[name]++
		}
	}

	var c Conflicts
	for name, n := range contexts {
		if n > 1 {
			c.Contexts = append(c.Contexts, name)
		}
	}
	for name, n := range clusters {
		if n > 1 {
			c.Clusters = append(c.Clusters, name)
		}
	}
	for name, n := range users {
		if n > 1 {
			c.Users = append(c.Users, name)
		}
	}
	return c
}
