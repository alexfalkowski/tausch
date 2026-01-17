# AGENTS.md

This repository is a small Go CLI + library that stubs command execution output based on a YAML config.

## Quick start (local)

### Prereqs

- Go toolchain: `go.mod` declares `go 1.25.0`.
- Some Make targets come from the `bin/` git submodule (see `.gitmodules`). If targets fail with missing files, initialize the submodule first:

```bash
make submodule
```

### Common workflow

```bash
make dep
make lint
make build
go test ./...
```

Notes:
- `go test ./...` expects the `tausch` binary to exist for `exec/exec_test.go`, so run `make build` (or `go build`) first.

## Essential commands

### Build

- Build the CLI binary in repo root:

```bash
make build
```

(`Makefile:6-7`)

### Dependencies

Provided via included makefiles from the `bin/` submodule (`Makefile:1-3` includes `bin/build/make/go.mak`):

```bash
make dep        # go mod download + tidy + vendor
make tidy
make vendor
make download
```

### Test

- Fast/local (requires `make build` first):

```bash
go test ./...
```

- CI-style specs (requires `gotestsum`):

```bash
make specs
```

(`bin/build/make/go.mak:62-64`)

### Lint / format

- Lint (runs field alignment + golangci-lint via the `bin/` submodule wrappers):

```bash
make lint
```

- Auto-fix (where supported):

```bash
make fix-lint
```

- Format:

```bash
make format
```

Lint configuration lives in `.golangci.yml`.

### Security

- Run govulncheck:

```bash
make sec
```

(`bin/build/make/go.mak:95-97`)

### Coverage

The CI pipeline generates coverage profiles under `test/reports/`.

```bash
make coverage
```

(`bin/build/make/go.mak:84-86`)

There are also helpers:

```bash
make html-coverage
make func-coverage
```

### Misc (only if you see them used)

These targets exist in `bin/build/make/go.mak` but require external tools:

- `make codecov-upload` (needs `codecovcli`)
- `make benchmark` / `make benchmark-pprof`
- `make create-diagram` (needs `goda`, `dot`)
- `make create-certs` (needs `mkcert`)
- `make analyse` (uses `gsa`)
- `make money` (uses `scc`)

## Project structure

- `main.go`: CLI entrypoint; calls `internal/cmd.Run(...)` and exits with returned status code.
- `internal/cmd/`: CLI orchestration (parse flags, load config, find command entry, write stdout/stderr).
- `internal/flag/`: argument parsing and config path resolution.
- `internal/config/`: YAML config decoding and command lookup.
- `internal/encoding/`: decodes `kind:data` values.
- `internal/io/`: writes decoded bytes to the requested writer.
- `exec/`: library wrapper around `os/exec` that runs the `tausch` binary transparently.
- `test/`: fixtures (example configs and recorded stdout/stderr) and `test/reports/` output directory.

## Runtime behavior and conventions

### CLI invocation shape

The CLI is designed to be invoked as:

```bash
tausch -config path/to/config.yml -- <your command as a string>
```

Internally, the “command name” is the args after `--` joined by spaces (`internal/flag/values.go:45-47`). That string must match the `name` field in the YAML config.

### Config file discovery

Config path is resolved in this order (`internal/flag/values.go:29-42`):

1. `-config <path>`
2. `TAUSCH_CONFIG` environment variable
3. `$HOME/.config/tausch/config.yml`

### Encoding format for stdout/stderr

`stdout` / `stderr` values in config are strings with a `kind:data` prefix (`internal/encoding/encoding.go:13-25`):

- `text:<literal text>`
- `base64:<base64-encoded bytes>`
- `file:<path to file>`

If the prefix is unknown, decoding fails with `ErrKindNotFound`.

### Library (`exec` package)

The `exec` package returns an `*exec.Cmd` that actually executes the `tausch` binary:

- It resolves the binary via `exec.LookPath("tausch")`, and falls back to `TAUSCH_PATH` if not found (`exec/exec.go:15-22`).
- It prefixes arguments with `--` before the real command (`exec/exec.go:24-26`).

Tests in `exec/exec_test.go` rely on the `tausch` binary being present (either via `PATH` including the repo root or via `TAUSCH_PATH`).

## Testing patterns

- Uses `github.com/stretchr/testify/require` assertions.
- Tests are table-driven in places (e.g., encoding/io decode/write tests).
- Fixture files live under `test/configs/`, `test/stdout/`, and `test/stderr/`.

## Style / formatting

- `.editorconfig` specifies:
  - Go files: tab indentation.
  - General files: 2-space indentation, LF endings, trim trailing whitespace.
- Linting is enforced via `.golangci.yml`.

## CI / release

- CircleCI config: `.circleci/config.yml` runs (in order) `make dep`, `make lint`, `make sec`, `make build`, `make specs`, `make coverage`.
- GoReleaser config: `.goreleaser.yml`.

## Gotchas

- **Submodule dependency**: Many `make` targets (lint/specs/coverage/sec/clean/etc.) are defined in `bin/build/make/*.mak` and call scripts under `bin/`. Ensure `bin/` submodule is initialized.
- **Tests require built binary**: `go test ./...` will fail unless `tausch` is built in the repo root (run `make build` first).
- **Command name matching**: command lookup is string-based (joined args after `--`); minor spacing differences will cause `command not found` errors.
