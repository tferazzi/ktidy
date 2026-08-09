# ktidy merge

Merge multiple kubeconfigs into one and print to stdout. Warns on duplicate context/cluster/user names.

```bash
ktidy merge work.yaml personal.yaml > ~/.kube/config
ktidy merge -p cluster-a.yaml cluster-b.yaml > merged.yaml  # keep external credential paths
```
