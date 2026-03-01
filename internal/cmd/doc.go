// Package cmd contains the CLI orchestration for tausch.
//
// This package is internal to the tausch module and is used by the main
// entrypoint to implement the `tausch` command's behavior.
//
// # Responsibilities
//
// The primary responsibility of this package is to:
//
//   - Parse CLI flags and remaining arguments (delegated to internal/flag).
//   - Resolve the tausch YAML config path (delegated to internal/flag).
//   - Decode the YAML configuration (delegated to internal/config).
//   - Look up the command entry by name (delegated to internal/config).
//   - Decode and write configured stdout/stderr payloads (delegated to internal/io).
//   - Return an exit status code that the CLI can use as its process exit code.
//
// # Command name matching
//
// The command "name" is derived from the arguments after flag parsing by joining
// them with spaces. That string must match the `name` field of a command in the
// configuration file.
//
// For example, if the user runs:
//
//	tausch -- go version
//
// then the derived command name is:
//
//	"go version"
//
// # Exit codes
//
// [Run] returns an integer exit code and an error.
//
// In normal operation:
//
//   - If the configured command's `stdout` value is present and successfully
//     written, [Run] returns exit code 0.
//   - Otherwise, it attempts to write the configured `stderr` value and returns
//     exit code 1.
//
// Errors are returned for failures such as flag parsing problems, missing or
// unreadable configuration, YAML decoding errors, command lookup failures, or
// decode/write errors while emitting stdout/stderr.
package cmd
