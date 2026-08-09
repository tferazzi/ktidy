# ktidy remove

Remove one or more contexts (and their referenced cluster + user) from the active kubeconfig.

```bash
ktidy remove old-cluster
ktidy remove staging dev old-cluster            # remove multiple at once
ktidy remove --dry-run prod                     # preview without writing
```
