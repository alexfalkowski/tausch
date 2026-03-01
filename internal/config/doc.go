// Package config provides decoding and lookup logic for tausch YAML configuration.
//
// This package is internal to the tausch module. It defines the in-memory
// representation of the YAML config file and helpers to decode it and find the
// configured behavior for a given command name.
//
// # YAML shape
//
// The configuration is expected to be a YAML document with a top-level `cmds`
// list. Each entry represents a command that tausch can stub.
//
// Example:
//
//	cmds:
//	  - name: "go version"
//	    stdout: "text:go version go1.25.0 darwin/arm64\n"
//	    stderr: ""
//
// # Stdout/stderr payload encoding
//
// The `stdout` and `stderr` fields are strings that are interpreted elsewhere
// (see internal/encoding and internal/io). They use a `kind:data` format such as:
//
//   - text:<literal text>
//   - base64:<base64-encoded bytes>
//   - file:<path to a file whose bytes should be written>
//
// This package treats those fields as opaque strings; it does not decode or
// validate the `kind:data` format.
//
// # Command lookup
//
// Commands are identified by their `name` field. The name is matched exactly
// against the command name derived by the CLI (typically the arguments after
// `--` joined with spaces).
//
// If a command cannot be found, [Config.GetCommand] returns [ErrCommandNotFound].
package config
