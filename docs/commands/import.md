# ktidy import

Merge kubeconfig(s) into the active kubeconfig (`$KUBECONFIG` or `~/.kube/config`). Skips conflicting entries unless `--force` is given.

```bash
ktidy import --save new-cluster.yaml           # write back to ~/.kube/config
ktidy import --save --force updated.yaml        # overwrite conflicting entries
ktidy import --dry-run new-cluster.yaml         # preview without writing
aws eks update-kubeconfig ... | ktidy import --stdin --save
```
