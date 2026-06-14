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

## ✨ Inspiration

I have taken these ideas from using tools from my past, such as:

- [vcr](https://github.com/vcr/vcr)
- [go-vcr](https://github.com/dnaeon/go-vcr)

Thank you for creating them.

> [!NOTE]
> One way to expand this tool in the future is to also run this once and record the outputs, if the need arises it will be added.

## ⚡ Quick Start

Install the CLI:

```bash
go install github.com/alexfalkowski/tausch@latest
```

Create a config:

```yaml
cmds:
  - name: go version
    stdout: "text:go version go1.26.0 darwin/arm64"
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

The configuration is a YAML document with a top-level `cmds` list. Each command entry has a `name` and either `stdout` or `stderr`.

The `stdout` and `stderr` values use a `kind:data` format:

```txt
text:This is awesome
base64:VGhpcyBpcyBhd2Vzb21l
file:path
```

The supported kinds are:

- `text:<literal text>` writes the text bytes as-is.
- `base64:<base64-encoded bytes>` decodes standard base64 and writes the bytes.
- `file:<path>` reads the file at the path and writes its bytes. Relative paths are resolved from the current working directory of the `tausch` process, not from the config file location.

The configuration can look like:

```yaml
cmds:
  - name: go version
    stdout: file:test/stdout/go_version.txt
  - name: go bob
    stderr: file:test/stderr/go_bob.txt
```

> [!IMPORTANT]
> Configure only one output stream per command. A command with `stdout` is treated as successful and exits `0`; a command with `stderr` and no `stdout` is treated as failing and exits `1`. Configuring both `stdout` and `stderr` for one command is rejected.

## 🎙️ Capture

To capture the combined `stdout` and `stderr` of a command, you can run:

```bash
command &> path
```

Examples:

```bash
go version &> test/stdout/go_version.txt
go bob &> test/stderr/go_bob.txt
```

> [!CAUTION]
> `&>` combines both streams. If your stub needs to prove that output came from only one stream, capture stdout with `>` or stderr with `2>` instead.

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

Build the CLI from the repository root:

```bash
make build
```

Run the tests after the binary exists:

```bash
go test ./...
```
