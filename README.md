# ktidy

A fast, safe kubeconfig manager. Merge, split, import, export, rename and group kubeconfig files — with conflict detection and `--dry-run` support throughout.

## Install

tbd.

## Commands

| Command                             | Description                                               |
| ----------------------------------- | --------------------------------------------------------- |
| [`merge`](docs/commands/merge.md)   | Merge multiple kubeconfigs into one                       |
| [`import`](docs/commands/export.md) | Merge kubeconfig(s) into the active kubeconfig            |
| [`export`](docs/commands/export.md) | Extract one or more contexts into a minimal kubeconfig    |
| [`split`](docs/commands/split.md)   | Export a context and optionally remove it from the source |
| [`remove`](docs/commands/remove.md) | Remove one or more contexts from the active kubeconfig    |
| [`rename`](docs/commands/rename.md) | Rename a context key                                      |
| [`group`](docs/commands/group.md)   | Organise kubeconfigs into named groups, merged at runtime |

## Global flags

| Flag        | Description                                       |
| ----------- | ------------------------------------------------- |
| `--dry-run` | Print what would happen without writing any files |
| `--version` | Print the version                                 |

## Shell completions

```bash
# zsh
ktidy completion zsh > $(brew --prefix)/share/zsh/site-functions/_ktidy

# bash
ktidy completion bash > /etc/bash_completion.d/ktidy
```

## Development

```bash
make build     # build ./bin/ktidy
make test      # run tests
make lint      # go vet
make snapshot  # local release build via goreleaser (no tag required)
```

## Releasing

1. Tag the commit: `git tag v0.x.y && git push --tags`
2. CI picks up the tag and runs `make release` via GoReleaser
3. Binaries for `linux/darwin/windows` and `amd64/arm64` are published to the GitHub Release automatically

To preview the release build locally without publishing: `make snapshot`
