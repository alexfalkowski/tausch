# Tausch

It is common to want to try to test commands part of the [exec](https://pkg.go.dev/os/exec) package.

This tool allows you to still call commands though just stub them out.

## Configuration

The configuration is just a lost of commands and wether you would like to write to stdout or stderr.

Each command can be `text`, `file` or a `base64` text.

```json
{
  "cmds": [
    {
      "name": "text_stdout",
      "stdout": "text:test"
    },
    {
      "name": "text_stderr",
      "stderr": "text:test"
    },
    {
      "name": "base64_stdout",
      "stdout": "base64:dGVzdA=="
    },
    {
      "name": "file_stdout",
      "stdout": "file:test/configs/test.txt"
    }
  ]
}
```

## Usage

> [!TIP]
> It is best to run your command and save the outputs to a file.

```bash
tausch -config test/configs/config.json -- text_stderr
```

You just pass in a config and after the `--` you call your usual command. So basically you just prefix your command with `tausch`.
