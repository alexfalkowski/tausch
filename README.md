[![CircleCI](https://circleci.com/gh/alexfalkowski/tausch.svg?style=shield)](https://circleci.com/gh/alexfalkowski/tausch)
[![codecov](https://codecov.io/gh/alexfalkowski/tausch/graph/badge.svg?token=AGP01JOTM0)](https://codecov.io/gh/alexfalkowski/tausch)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/tausch)](https://goreportcard.com/report/github.com/alexfalkowski/tausch)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/tausch.svg)](https://pkg.go.dev/github.com/alexfalkowski/tausch)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)

# 🔁 Tausch

It is common to test Go code that uses the [exec](https://pkg.go.dev/os/exec) package.

Tausch lets that code keep calling commands while replacing the command output with configured stubs.

## 📚 Background

Writing tools in the Unix/Linux world is a composition of many other tools. This is, in a lot of ways, the [unix philosophy](https://en.wikipedia.org/wiki/Unix_philosophy).

Though as we start to write these scripts/tools we quickly realise that verifying they work is not easy. You might ask yourself how can we test these?

If you have done some [test-driven development](https://en.wikipedia.org/wiki/Test-driven_development), you might be wondering the same thing?

Now the world of writing testable scripts has come a long way from when we first learned about [shell scripts](https://en.wikipedia.org/wiki/Shell_script).

For those interested, check out the following:

- [Bats-core: Bash Automated Testing System](https://github.com/bats-core/bats-core)
- [shUnit2](https://github.com/kward/shunit2)
- [Aruba](https://github.com/cucumber/aruba)
- [Advanced Bash-Scripting Guide](https://tldp.org/LDP/abs/html/index.html)

Though as your codebase starts getting bigger, you start to question was this scripting language the right choice? This is where you will get other solutions, such as:

- [script](https://github.com/bitfield/script)
- [zx](https://github.com/google/zx)
- [abs](https://github.com/abs-lang/abs)
- [wren](https://github.com/wren-lang/wren)

All these projects are fine, though you might have some hard requirements or just preferences and need to use [go](https://go.dev/). If this is you, then this project might help.

## 🧭 Drift

As this tool just [stubs](http://xunitpatterns.com/Test%20Stub.html) the output of a command, how do we make sure that we stay compatible with what we record?

Well this is a tough problem, though a problem regardless of this tool. As the tools you use will change in subtle (or not) and interesting ways.

This is known as [dependency hell](https://en.wikipedia.org/wiki/Dependency_hell). Upgrading dependencies need to be verified for compatibility. Practices like [semantic Versioning](https://semver.org/) and to [pin exact dependency versions](https://betterdev.blog/pin-exact-dependency-versions/), will likely help or you might say that is [wishful thinking](https://en.wikipedia.org/wiki/Wishful_thinking).

So managing the outputs needs a careful process. One ways it to have a [runbook](https://en.wikipedia.org/wiki/Runbook) or find a way run it periodically to record the output.

Though as you might realise it does not matter if you manage to verify/tests lots of combinations, you will miss something. This is why making sure that whatever you build you need to have [observability](https://en.wikipedia.org/wiki/Observability_(software)) as a first class citizen.

## ❓ Why?

So you might be asking yourself, why should I use this solution?

Some of these reasons I have encountered:

- Commands can take a long time to run.
- Dependency setup can be costly.
- Simulating failure can be hard.

As this is a single binary and ties into the already defined [cmd](https://pkg.go.dev/os/exec#Cmd) type in the defined stdlib.

Of course you might not want another dependency, if that is the case then just copy the [code](https://github.com/alexfalkowski/tausch/blob/master/exec/exec.go).

## ⚡ Quick Start

Tausch requires a Go toolchain compatible with the module's `go.mod` directive.

Install the CLI:

```bash
go install github.com/alexfalkowski/tausch@latest
```

Create a config:

```yaml
cmds:
  - name: go version
    stdout: "text:stubbed go version"
  - name: go bob
    stderr: "text:go bob: unknown command"
```

Run a command through Tausch:

```bash
tausch -config config.yml -- go version
```

> [!TIP]
> Set `TAUSCH_CONFIG` in tests when many commands should share the same config path.

## ⚙️ Configuration

The configuration is a YAML document with a top-level `cmds` list. Each command entry has a `name` and may set `stdout`, `stderr`, or `exit_code`.

The `stdout` and `stderr` values use a `kind:data` format:

```txt
text:This is awesome
base64:VGhpcyBpcyBhd2Vzb21l
file:path
```

The supported kinds are:

- `text:<literal text>` writes the text bytes as-is.
- `base64:<base64-encoded bytes>` decodes standard base64 and writes the bytes.
- `file:<path>` reads the file at the path and writes its bytes. Relative paths are resolved from the directory containing the config file.

The configuration can look like:

```yaml
cmds:
  - name: go version
    stdout: file:test/stdout/go_version.txt
  - name: go bob
    stderr: file:test/stderr/go_bob.txt
    exit_code: 127
  - name: grep needle file.txt
    exit_code: 1
```

> [!IMPORTANT]
> Configure at most one non-empty encoded output stream for each stub. By default, a non-empty `stdout` value exits `0`; otherwise, Tausch falls back to `stderr` and exits `1`. If both streams are empty or omitted and `exit_code` is not set, the command exits `1` with no configured output. Set `exit_code` to override the final status after any configured output is written successfully; valid values are `0` through `255`. To stub a successful command that writes zero bytes, use `stdout: "text:"`. Configuring both `stdout` and `stderr` for one command is rejected.

Payload values and `file:` paths are decoded only when the matching command is invoked, not when the YAML file is loaded. Unknown kinds or payloads without the `kind:data` separator make that invocation exit `1` with the decode error.

## 🎙️ Capture

Capture real command output with shell redirection, then reference the recorded files from the config.

For stdout:

```bash
go version > test/stdout/go_version.txt
```

For stderr:

```bash
go bob 2> test/stderr/go_bob.txt
echo $? # copy this into exit_code when needed
```

For a command where you want both streams in one file:

```bash
command &> path
```

Then configure the recorded output:

```yaml
cmds:
  - name: go version
    stdout: file:test/stdout/go_version.txt
  - name: go bob
    stderr: file:test/stderr/go_bob.txt
    exit_code: 127
```

To refresh a fixture, rerun the same redirect command:

```bash
go version > test/stdout/go_version.txt
```

> [!CAUTION]
> `&>` combines stdout and stderr. Use `>` or `2>` when the stub should prove which stream produced the output.

## 🚀 Usage

There are multiple ways you can use this.

### 🌱 Environment

The executable looks for configuration in a few places.

#### 🧩 Executable

The executable reads the config path in this order:

- `-config` - argument with a path.
- `TAUSCH_CONFIG` - from an env variable.
- The platform user config directory with `tausch/config.yml` appended - default config location. On Unix-like systems this is commonly `$HOME/.config/tausch/config.yml`, but the exact base directory comes from Go's `os.UserConfigDir`.

#### 🧪 exec.CommandContext

Using the library looks for the executable in this order:

- `$PATH` - finds the `tausch` executable provided in the path.
- `TAUSCH_PATH` - path of the binary if `tausch` is not found on `PATH`.

PATH lookup uses Go's `os/exec.LookPath("tausch")`. If Go rejects a relative current-directory match such as `PATH=.`, the wrapper falls back to `TAUSCH_PATH`; use an absolute directory on `PATH` or set `TAUSCH_PATH` in tests.

### 💻 Command

Pass a config and, after `--`, call your usual command. The command tokens after `--` are joined with spaces and matched against a config entry's `name`.

> [!WARNING]
> Command matching is exact and case-sensitive. For example, `tausch -- go version` matches `name: go version`, but not `name: Go Version` or `name: go version extra`.

#### 🆘 Help

Print command usage with any of the standard Go flag help forms:

```bash
tausch -h
tausch -help
tausch --help
```

Help output includes the invocation shape, config path resolution order, and available flags. These help invocations exit with status `1`.

#### ✅ Example for stdout

```bash
tausch -config test/configs/config.yml -- go version
```

#### ❌ Example for stderr

```bash
tausch -config test/configs/config.yml -- go bob
```

To verify it caused an error:

```bash
echo $? # 1
```

### 📦 Library

Add the module dependency in the Go module that uses the wrapper:

```bash
go get github.com/alexfalkowski/tausch@latest
```

In your code you would use it just like you would the [exec](https://pkg.go.dev/os/exec).

```go
import (
  "context"

  "github.com/alexfalkowski/tausch/exec"
)

cmd := exec.CommandContext(context.Background(), "go", "version")
_ = cmd.Run()
```

The wrapper above invokes the `tausch` binary with `-- go version`, so the command name must be present in the active YAML config.

For library use, configure Tausch with `TAUSCH_CONFIG` or the default config location before running the command. `CommandContext` arguments are only the target command tokens; passing `-config` there would make it part of the stubbed command name instead of a Tausch CLI flag.

## 🛠️ Development

For a fresh checkout, initialize the shared build tooling and module dependencies first:

```bash
git submodule update --init
make dep
```

Build the CLI from the repository root:

```bash
make build
```

Run the repository checks after the binary exists:

```bash
make specs
```

Run the benchmark aggregate:

```bash
make benchmarks
```

Run one benchmark package:

```bash
make benchmark package=internal/io
```

Run count-bounded smoke fuzzing for all fuzz targets:

```bash
make fuzzes
```

Run one fuzz target:

```bash
make fuzz package=internal/io name=FuzzWriteText
```

The repository default is `fuzztime=1000x`, so fuzzing is bounded by run count
instead of wall-clock duration. Override it when needed:

```bash
make fuzz package=internal/io name=FuzzWriteText fuzztime=250x
```
