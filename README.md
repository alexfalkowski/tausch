[![CircleCI](https://circleci.com/gh/alexfalkowski/tausch.svg?style=shield)](https://circleci.com/gh/alexfalkowski/tausch)
[![codecov](https://codecov.io/gh/alexfalkowski/tausch/graph/badge.svg?token=AGP01JOTM0)](https://codecov.io/gh/alexfalkowski/tausch)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/tausch)](https://goreportcard.com/report/github.com/alexfalkowski/tausch)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/tausch.svg)](https://pkg.go.dev/github.com/alexfalkowski/tausch)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)

# Tausch

It is common to want to try to test commands part of the [exec](https://pkg.go.dev/os/exec) package.

This tool allows you to still call commands though just stub them out.

## Configuration

The configuration is just a lost of commands and wether you would like to write to stdout or stderr.

Each command can be `text`, `file` or a `base64` text.

```yaml
cmds:
- name: text_stdout
  stdout: text:test
- name: text_stderr
  stderr: text:test
- name: base64_stdout
  stdout: base64:dGVzdA==
- name: file_stdout
  stdout: file:test/configs/test.txt
```

## Usage

> [!TIP]
> It is best to run your command and save the outputs to a file.

```bash
tausch -config test/configs/config.yaml -- text_stderr
```

You just pass in a config and after the `--` you call your usual command. So basically you just prefix your command with `tausch`.
