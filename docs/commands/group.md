# ktidy group

Organise kubeconfigs into named directories under `~/.kube/groups/`. Every file in a group directory is a separate kubeconfig; they are merged at runtime via `KUBECONFIG`.

```text
~/.kube/groups/
  aws/
    eks-prod.yaml
    eks-staging.yaml
  gcp/
    gke-europe.yaml
```

```bash
ktidy group create aws
ktidy group add aws ~/Downloads/eks-prod.yaml
ktidy group add --move aws ~/Downloads/eks-staging.yaml
ktidy group list
ktidy group activate              # all groups
ktidy group activate aws gcp      # named groups only
ktidy group remove old-group
```

Provider workflow — write directly into the group directory, no `ktidy import` needed:

```bash
gcloud container clusters get-credentials my-cluster \
  --kubeconfig ~/.kube/groups/gcp/my-cluster.yaml

aws eks update-kubeconfig --name my-cluster \
  --kubeconfig ~/.kube/groups/aws/my-cluster.yaml
```

Override the groups root with `KTIDY_GROUPS_DIR` or `--groups-dir`.

## Shell setup

For `ktidy group` to work across new shells, add this once to `~/.zshrc` or `~/.bashrc`:

```bash
eval "$(ktidy group activate)"
```
