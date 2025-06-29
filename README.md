[![CircleCI](https://circleci.com/gh/alexfalkowski/tausch.svg?style=shield)](https://circleci.com/gh/alexfalkowski/tausch)
[![codecov](https://codecov.io/gh/alexfalkowski/tausch/graph/badge.svg?token=AGP01JOTM0)](https://codecov.io/gh/alexfalkowski/tausch)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/tausch)](https://goreportcard.com/report/github.com/alexfalkowski/tausch)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/tausch.svg)](https://pkg.go.dev/github.com/alexfalkowski/tausch)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)

# Tausch

It is common to want to try to test commands part of the [exec](https://pkg.go.dev/os/exec) package.

This tool allows you to still call commands though just stub them out.

## Background

Writing tools in the Unix/Linux world is a composition of many other tools. This is in a lot of ways is the [Unix philosophy](https://en.wikipedia.org/wiki/Unix_philosophy)

Though as we start to write these scripts/tools we quickly realise that verifying they work is not easy. You might ask yourself how can we test these?

If you have done some [Test-driven development](https://en.wikipedia.org/wiki/Test-driven_development), you might be wondering the same thing?

Now the world of writing testable scripts has come a long way from when we first learned about [shell scripts](https://en.wikipedia.org/wiki/Shell_script).

For those interested, check out the following:

- [Bats-core: Bash Automated Testing System](https://github.com/bats-core/bats-core)
- [shUnit2](https://github.com/kward/shunit2)
- [Aruba](https://github.com/cucumber/aruba)
- [Advanced Bash-Scripting Guide](https://tldp.org/LDP/abs/html/index.html)

Though as your codebase starts getting bigger, you start to question was this scripting language the right choice? This is where you will get other solutions, such as:

- [zx](https://github.com/google/zx)
- [abs](https://github.com/abs-lang/abs)
- [wren](https://github.com/wren-lang/wren)

All these projects are fine, though you might have some hard requirements or just preferences and need to use [Go](https://go.dev/). If this is you, then this project might help.

## Drift

As this tool just fakes the output of a command, how do we make sure that we stay compatible with what we record?

Well this is a tough problem, though a problem regardless of this tool. As the tools you use will change in subtle (or not) and interesting ways.

This is known as [dependency management problem](https://en.wikipedia.org/wiki/Dependency_hell). Upgrading dependencies need to be verified for compatibility. Practices like [Semantic Versioning](https://semver.org/) and [Pin exact dependency versions](https://betterdev.blog/pin-exact-dependency-versions/), might help ([wishful thinking](https://en.wikipedia.org/wiki/Wishful_thinking)).

So when updating the outputs when using this tool have a guided [Runbook](https://en.wikipedia.org/wiki/Runbook) on how to manage these (or automate if you can).

Though as you might realise it does not matter if you manage to verify/tests lots of combinations, you will miss something. This is why making sure that whatever you build you need to have [Observability](https://en.wikipedia.org/wiki/Observability_(software)) as a first class citizen.

## Configuration

The configuration is just a list of `cmds` and wether you would like to write to `stdout` or `stderr`.

Each `cmd` can be `text`, `file` or a `base64` text, separated by a `: (colon)`.

Examples:

```txt
text:This is awesome
base64:VGhpcyBpcyBhd2Vzb21l
file:path
```

The configuration would look like:

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

Examples:

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
