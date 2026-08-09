# ktidy

A Go implementation to manage kubeconfig files — merge, split, import, export and group kubeconfig files.

## Project layout

```text
cmd/            # one file per subcommand; no business logic here
internal/
  kubeconfig/   # all clientcmd wrappers and kubeconfig mutation logic
main.go         # entry point — only wires version and calls cmd.Execute
```

All kubeconfig reads/writes go through `internal/kubeconfig`. `cmd/` files only parse flags, call internal functions, and print results.

## Go style

Follow the [Google Go Style Guide](https://google.github.io/styleguide/go/). Key rules that apply here:

- Package names: lowercase, single word, no underscores (`kubeconfig`, not `kube_config`)
- Error strings: lowercase, no punctuation (`"context not found"`, not `"Context not found."`)
- Return errors; don't log or print inside `internal/` — let `cmd/` decide how to surface them
- No `init()` for anything other than registering cobra subcommands and flags
- Prefer `errors.New` / `fmt.Errorf("...: %w", err)` over custom error types unless the caller needs to type-assert
- Table-driven tests with `t.Run` for any non-trivial logic in `internal/`

## Cobra conventions

- Each subcommand lives in its own file: `cmd/<name>.go`
- Register with `rootCmd.AddCommand(...)` inside an `init()` at the bottom of the file
- Use `RunE` (returns error) not `Run` — so errors propagate to `cobra.CheckErr` in `Execute`
- Persistent flags on `rootCmd` only (e.g. `--groups-dir`); local flags on the subcommand
- Use `MarkFlagsRequiredTogether` / `MarkFlagsMutuallyExclusive` instead of manual validation where applicable
- Dynamic completions via `ValidArgsFunction` for any arg that takes a context or group name; return `cobra.ShellCompDirectiveNoFileComp` when file completion is not wanted
- `Use` field format: `"<verb> [flags] <ARGS>"` — positional args in UPPER_SNAKE_CASE

## Kubeconfig conventions

- Always load via `internal/kubeconfig.Load(path)` — never call `clientcmd` directly from `cmd/`
- Always write via `internal/kubeconfig.Write(cfg, path)` — this is where `--dry-run` will be enforced (issue #10)
- Conflict detection happens before any write; warn to stderr, never silently overwrite
- Default config path resolution order: `-k` flag → `KUBECONFIG` env → `~/.kube/config`

## Agent skills

### Issue tracker

Issues live in GitHub Issues for this repo. See `docs/agents/issue-tracker.md`.

### Triage labels

Default five-role vocabulary (needs-triage → wontfix). See `docs/agents/triage-labels.md`.

### Domain docs

Single-context layout — `CONTEXT.md` + `docs/adr/` at repo root. See `docs/agents/domain.md`.
