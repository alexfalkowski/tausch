// Package exec provides a small wrapper around the standard library's os/exec
// that transparently invokes the tausch CLI.
//
// # Overview
//
// The goal of this package is to let you construct an *exec.Cmd exactly like
// you would with [os/exec], but have the command routed through the `tausch`
// binary so that stdout/stderr and exit codes can be stubbed based on a
// tausch YAML configuration.
//
// Concretely, [CommandContext] creates an [os/exec.Cmd] whose executable is
// the `tausch` binary and whose arguments are prefixed with `--`, followed by
// the command you want to run.
//
// # Executable resolution
//
// The `tausch` executable path is resolved as follows:
//
//  1. If `tausch` can be found on PATH via [os/exec.LookPath], that path is used.
//  2. Otherwise, the value of the `TAUSCH_PATH` environment variable is used.
//
// If neither produces a usable path, the returned command will fail when run
// (for example with an “executable file not found” error).
//
// # Configuration
//
// The tausch CLI discovers its YAML config path via `-config`, `TAUSCH_CONFIG`,
// or a default location. This package does not set or interpret configuration
// itself; you are expected to configure tausch via environment variables or
// CLI flags as appropriate for your application/test.
//
// # Example
//
//	cmd := exec.CommandContext(ctx, "go", "version")
//	out, err := cmd.CombinedOutput()
//
// This will execute something equivalent to:
//
//	tausch -- go version
//
// using the same semantics for stdin/stdout/stderr handling as [os/exec.Cmd].
package exec
