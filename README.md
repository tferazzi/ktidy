# ktidy

A fast, safe kubeconfig manager. Merge, split, import, export, rename, and group kubeconfig files — with conflict detection and `--dry-run` support throughout.

## Install

```bash
# Homebrew
brew install tferazzi/ktidy/ktidy

# Direct download (macOS arm64)
curl -Lo ktidy https://github.com/tferazzi/ktidy/releases/latest/download/ktidy_darwin_arm64
chmod +x ktidy && mv ktidy /usr/local/bin
```

## Shell setup

For `ktidy group` to work across new shells, add this once to `~/.zshrc` or `~/.bashrc`:

```bash
eval "$(ktidy group activate)"
```

## Commands

### `ktidy merge`

Merge multiple kubeconfigs into one and print to stdout. Warns on duplicate context/cluster/user names.

```bash
ktidy merge work.yaml personal.yaml > ~/.kube/config
ktidy merge -p cluster-a.yaml cluster-b.yaml > merged.yaml  # keep external credential paths
```

### `ktidy import`

Merge kubeconfig(s) into the active kubeconfig (`$KUBECONFIG` or `~/.kube/config`). Skips conflicting entries unless `--force` is given.

```bash
ktidy import --save new-cluster.yaml           # write back to ~/.kube/config
ktidy import --save --force updated.yaml        # overwrite conflicting entries
ktidy import --dry-run new-cluster.yaml         # preview without writing
aws eks update-kubeconfig ... | ktidy import --stdin --save
```

### `ktidy export`

Extract one or more contexts into a minimal kubeconfig and print to stdout.

```bash
ktidy export prod > prod.yaml
ktidy export -k ~/other.yaml staging prod > two-contexts.yaml
```

### `ktidy split`

Like `export`, but `--remove` also deletes the context from the source file.

```bash
ktidy split prod > prod.yaml                    # export only
ktidy split --remove prod > prod.yaml           # export and remove from source
ktidy split --remove --dry-run prod             # preview the removal
```

### `ktidy remove`

Remove one or more contexts (and their referenced cluster + user) from the active kubeconfig.

```bash
ktidy remove old-cluster
ktidy remove staging dev old-cluster            # remove multiple at once
ktidy remove --dry-run prod                     # preview without writing
```

### `ktidy rename`

Rename a context key. Updates `current-context` if it matches.

```bash
ktidy rename arn:aws:eks:eu-west-1:123:cluster/prod prod
ktidy rename --dry-run old-name new-name
```

### `ktidy group`

Organise kubeconfigs into named directories under `~/.kube/groups/`. Every file in a group directory is a separate kubeconfig; they are merged at runtime via `KUBECONFIG`.

```
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

## Global flags

| Flag | Description |
|------|-------------|
| `--dry-run` | Print what would happen without writing any files |
| `--version` | Print the version |

## Shell completions

```bash
# zsh
ktidy completion zsh > $(brew --prefix)/share/zsh/site-functions/_ktidy

# bash
ktidy completion bash > /etc/bash_completion.d/ktidy
```
