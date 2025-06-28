[![CircleCI](https://circleci.com/gh/alexfalkowski/tausch.svg?style=shield)](https://circleci.com/gh/alexfalkowski/tausch)
[![codecov](https://codecov.io/gh/alexfalkowski/tausch/graph/badge.svg?token=AGP01JOTM0)](https://codecov.io/gh/alexfalkowski/tausch)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/tausch)](https://goreportcard.com/report/github.com/alexfalkowski/tausch)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/tausch.svg)](https://pkg.go.dev/github.com/alexfalkowski/tausch)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)

# Tausch

It is common to want to try to test commands part of the [exec](https://pkg.go.dev/os/exec) package.

This tool allows you to still call commands though just stub them out.

## Configuration

The configuration is just a list of `cmds` and wether you would like to write to `stdout` or `stderr`.

Each `cmd` can be `text`, `file` or a `base64` text.

Example:

```yaml
cmds:
- name: go version
  stdout: file:test/stdout/go_version.txt
- name: go bob
  stderr: file:test/stderr/go_bob.txt
```

## Capture

To capture the `stdout` or `stderr` of the command, you can run the following:

```bash
command &> path
```

### Examples

```bash
go version &> test/stdout/go_version.txt
go bob &> test/stderr/go_bob.txt
```

## Usage

There are multiple ways you can use this.

### Environment

The executable looks for configuration in a few places.

#### Executable

The executable will read the config from the following places:

- `-config` - argument with a path.
- `TAUSCH_CONFIG` - from an env variable.
- `$HOME/.config/tausch.yml` - The config can be placed a well known config folder.

#### exec.CommandContext

Using the library will look for the executable in the following places:

- `TAUSCH_PATH` - the path of the binary.
- `$PATH`: - finds the executable provided in the path.

### Command

You just pass in a config and after the `--` you call your usual command. So basically you just prefix your command with `tausch`.

#### Example for `stdout`

```bash
tausch -config test/configs/config.yml -- go version
```

#### Example for `stderr`

```bash
tausch -config test/configs/config.yml -- go bob
```

To verify it caused and error:

```bash
echo $? => 1
```

### Library

In your code you would use it just like you would the [exec](https://pkg.go.dev/os/exec).

```go
import (
  "context"

  "github.com/alexfalkowski/tausch/exec"
)

cmd := exec.Command("go", "version")
_ = cmd.Run()

cmd := exec.CommandContext(context.Background(), "go", "version")
_ = cmd.Run()
```
