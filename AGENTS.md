# AGENTS.md

This repository is a small Go CLI + library that stubs command execution output based on a YAML config.

## Shared guidance

Use `bin/AGENTS.md` for shared skills and cross-repository defaults.

## Quick start (local)

### Prereqs

- Go toolchain details live in `go.mod`.
- Some Make targets come from the `bin/` git submodule (see `.gitmodules`). If
  targets fail with missing files, initialize the submodule first. Use
  `make submodule` once the shared checkout is present; see `bin/AGENTS.md` for
  fresh-clone bootstrap details.

```bash
make submodule
```

### Common workflow

```bash
make dep
make lint
make build
make specs
```

Notes:
- The specs expect the `tausch` binary to exist for `exec/exec_test.go`, so run
  `make build` first when needed.
- Do not add a local `specs` target or `specs` prerequisite in this Makefile;
  `specs` is owned by the shared `bin/build/make/go.mak` fragment. CI encodes
  the required `make build` before `make specs` ordering.

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
make dep
make tidy
make vendor
make download
```

### Test

- Repository specs:

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

- Run the repository security target:

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

These targets exist in `bin/build/make/go.mak` but may require optional
external tools:

- `make codecov-upload`
- `make benchmark` / `make benchmark-pprof`
- `make create-diagram`
- `make create-certs`
- `make analyse`
- `make money`

## Project structure

- `main.go`: CLI entrypoint; calls `internal/cmd.Run(...)` and exits with returned status code.
- `internal/cmd/`: CLI orchestration (parse flags, load config, find command entry, write stdout/stderr).
- `internal/flag/`: argument parsing and config path resolution.
- `internal/config/`: YAML config decoding and command lookup.
- `internal/io/`: decodes `kind:data` values and writes bytes to the requested writer.
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
3. `os.UserConfigDir()/tausch/config.yml` (for example `$HOME/.config/tausch/config.yml` on many Unix-like systems)

### Encoding format for stdout/stderr

`stdout` / `stderr` values in config are strings with a `kind:data` prefix:

- `text:<literal text>`
- `base64:<base64-encoded bytes>`
- `file:<path to file>`

Relative `file:` paths are resolved from the directory containing the config file.

If the prefix is unknown, decoding fails with `ErrKindNotFound`.

### Library (`exec` package)

The `exec` package returns an `*exec.Cmd` that actually executes the `tausch` binary:

- It resolves the binary via `exec.LookPath("tausch")`, and falls back to `TAUSCH_PATH` if not found (`exec/exec.go:15-22`).
- It prefixes arguments with `--` before the real command (`exec/exec.go:24-26`).

Tests in `exec/exec_test.go` rely on the `tausch` binary being present (either via `PATH` including an absolute repo root or via `TAUSCH_PATH`). Do not rely on `PATH=.`; Go's `exec.LookPath` can reject current-directory matches and cause the wrapper to fall back to `TAUSCH_PATH`.

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

- CircleCI config: `.circleci/config.yml` wraps the core local validation sequence (`make dep`, `make lint`, `make sec`, `make build`, `make specs`, `make coverage`) with CI setup, cleanup, cache, and upload steps.
- GoReleaser config: `.goreleaser.yml`.

## Gotchas

- **Submodule dependency**: Many `make` targets (lint/specs/coverage/sec/clean/etc.) are defined in `bin/build/make/*.mak` and call scripts under `bin/`. Ensure `bin/` submodule is initialized.
- **Tests require built binary**: specs need the `tausch` binary in the repo
  root, so run `make build` first when needed. Do not override or extend the
  shared `specs` target locally; keep the build-before-specs ordering in
  workflow guidance and CI.
- **Command name matching**: command lookup is string-based (joined args after `--`); minor spacing differences will cause `command not found` errors.
